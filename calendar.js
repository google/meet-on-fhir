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

const {google} = require('googleapis');

exports.createEvent = function(client, encounterId, callback) {
	const start = new Date();
	const end = new Date(start.getTime() + 30 * 60 * 1000);
	const event = {
		summary: 'Hangouts Meet',
		start: {
			dateTime: start.toISOString(),
		},
		end: {
			dateTime: end.toISOString(),
		},
		conferenceData: {
			createRequest: {
				requestId: Math.random().toString(),
				conferenceSolutionKey: {
					type: "hangoutsMeet",
				},
			},
		},
	};

	const calendar = google.calendar({version: 'v3', auth: client});
	calendar.events.insert({
		calendarId: 'primary',
		conferenceDataVersion: 1,
		resource: event,
	}, (err, result) => {
    var link;
    if (result && result.data && result.data.hangoutLink) {
      link = result.data.hangoutLink;
    }
		callback(err, link);
	});
};
