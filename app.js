/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const calendar = require('./calendar.js');
const datastore = require('./datastore.js');
const user = require('./user.js');
const mllp = require('./mllp.js');
const fhir = require('./fhir.js');

const settings = require('./settings.json');

const express = require('express');
const session = require('cookie-session');
const e = require('express');

const app = express();
app.use(express.static('static'));
app.use('/fhirclient', express.static('node_modules/fhirclient/build/'));
app.use('/jquery', express.static('node_modules/jquery/dist/'));
app.use(express.urlencoded({extended: false}));
app.use(session({
	name: 'session',
	keys: [settings.sessionCookieSecret],
	maxAge: 7 * 60 * 60 * 1000,
}));

function error(response) {
  return function(err) {
    console.log(err);
    response.status(500).send(err);
  };
}

function debugLog(message) {
	if (settings.debugLogging) {
		console.log(message);
	}
}

app.get('/hangouts/:encounterId', (request, response) => {
	const key = datastore.key(['Encounter', request.params.encounterId]);
	datastore.get(key).then(entity => {
		if (entity) {
			debugLog('Patient encounter ' + request.params.encounterId + ' found URL ' + entity.Url);
			response.send({url: entity.Url});
		} else {
			response.send({});
		}
	}).catch(error(response));
});

app.post('/hangouts', (request, response) => {
	const encounterId = request.body.encounterId;
	const key = datastore.key(['Encounter', encounterId]);
	datastore.get(key).then(entity => {
		if (entity) {
			debugLog('Provider found existing encounter ' + request.body.encounterId + ' with URL ' + entity.Url);
			response.send({url: entity.Url});
			return;
		}

		user.withCredentials(request, response, client => {
			calendar.createEvent(client, encounterId, (err, url) => {
				if (err) {
					debugLog('ERROR: Provider calendar event create for encounter ' + request.body.encounterId + ' failed with error ' + err);
					response.status(500).send(err);
					return;
				}
				debugLog('Provider created calendar event for encounter ' + request.body.encounterId + ' with URL ' + url);
				const entity = { Url: url };
				datastore.set(key, entity).then(() => {
					response.send({url: url});
				});
			});
		});
	}).catch(error(response));
});

app.post('/reportEvent', async (request, response) => {
	const fhirUrl = request.body.fhirUrl;
	if (!fhirUrl) {
		response.status(400).send("missing fhirUrl");
		return;
	}
	if (!settings.authorizedFhirUrls.includes(fhirUrl)) {
		response.status(403).send(`unauthorized fhirUrl ${fhirUrl}`);
		return;
	}

	const fhirAccessToken = request.body.fhirAccessToken;
	if (!fhirAccessToken) {
		response.status(400).send("missing fhirAccessToken");
		return;
	}
  
	const type = request.body.type;
	if (type != 'patient_arrived' && type != 'practitioner_arrived') {
		response.status(400).send(`Invalid type ${type}`);
		return;
	}
  
	const encounterId = request.body.encounterId;
	if (!encounterId) {
		response.status(400).send('missing encounterId');
		return;
	}

	const patientId = request.body.patientId;
	if (!patientId) {
		response.status(400).send('missing patientId');
		return
	}

	const patientName = request.body.patientName;
	if (!patientName) {
		response.status(400).send('missing patientName');
		return
	}

	try {
		await fhir.checkFhirAuthorization(fhirUrl, fhirAccessToken, encounterId);
	} catch (err) {
		debugLog('fhir authentication check failed with err ' + err);
		response.status(403).send('fhir authentication check failed');
		return
	}
	
	const key = datastore.key(['Encounter', encounterId]);
	datastore.get(key).then(entity => {
		if (!entity) {
			response.status(400).send(`Meeting not found for encounter ${encounterId}`);
			return;
		}

		if (type == 'patient_arrived') {
			if (entity.patientArriveTime) {
				response.status(202).send(`patient has arrived before`);
				return;
			}
			debugLog('Updating patient arrive time for encounter ' + request.body.encounterId);
			entity.patientArriveTime = Date.now();
		}
		if (type == 'practitioner_arrived') {
			if (entity.practitionerArriveTime) {
				response.status(202).send(`practitioner has arrived before`);
				return;
			}
			debugLog('Updating practitioner arrive time for encounter ' + request.body.encounterId);
			entity.practitionerArriveTime = Date.now();
		}
		return datastore.update(key, entity);
	}).then(entity => {
		if (!entity || !settings.enableEHRWriteback) {
			response.status(200).send('EHR writeback is not needed or disabled.');
			return;
		}
		if (type == 'patient_arrived') {
			debugLog('Sending patient arrived notification to EHR for encounter ' + request.body.encounterId);
			//TODO: Send patient arrived message to EHR.
		}
		if (entity.practitionerArriveTime && entity.patientArriveTime) {
			debugLog('Sending appointment status change message to EHR for encounter ' + request.body.encounterId);
			return mllp.setAppointmentStatusArrived(encounterId, patientId, patientName);
		}
	}).then(() => {
		response.sendStatus(200);
	}).catch(error(response));
});

app.get('/authenticate', (request, response) => {
	user.authenticate(request, response);
});

app.get('/logout', (request, response) => {
	user.logout(request, response);
});

app.get('/settings', (request, response) => {
  response.send({'fhirClientId': settings.fhirClientId});
});

app.listen(process.env.PORT || 8080);
