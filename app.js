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

const settings = require('./settings.json');

const express = require('express');
const session = require('cookie-session');

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

app.get('/hangouts/:encounterId', (request, response) => {
	const key = datastore.key(['Encounter', request.params.encounterId]);
	datastore.get(key).then(entity => {
		if (entity) {
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
			response.send({url: entity.Url});
			return;
		}

		user.withCredentials(request, response, client => {
			calendar.createEvent(client, encounterId, (err, url) => {
				if (err) {
					response.status(500).send(err);
					return;
				}
				const entity = { Url: url };
				datastore.set(key, entity).then(() => {
					response.send({url: url});
				});
			});
		});
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
