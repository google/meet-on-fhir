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

const settings = require('./settings.json');

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

  withCalendarId(calendar, (err, id) => {
    if (err) {
      callback(err, null);
      return;
    }

    calendar.events.insert({
      calendarId: id,
      conferenceDataVersion: 1,
      resource: event,
    }, (err, result) => {
      var link;
      if (result && result.data && result.data.hangoutLink) {
        link = result.data.hangoutLink;
      }
      callback(err, link);
    });
  });
};

function withCalendarId(calendar, callback) {
  if (!settings.calendar || settings.calendar == 'primary') {
    callback(null, 'primary');
    return;
  }

  calendar.calendarList.list({ minAccessRole: 'owner' }, (err, result) => {
    if (result && result.data) {
      const items = result.data.items;
      for (var i = 0; i < items.length; i++) {
        if (items[i].summary == settings.calendar) {
          callback(err, items[i].id);
          return;
        }
      }

      calendar.calendars.insert({ requestBody: { summary: settings.calendar } }, (err, result) => {
        var id;
        if (result && result.data && result.data.id) {
          id = result.data.id;
        }
        callback(err, id);
      });
      return;
    }

    callback(err, null);
  })
}
