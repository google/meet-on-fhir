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

const datastore = require('./datastore.js');

const settings = require('./settings.json');

const crypto = require('crypto');
const {google} = require('googleapis');

function newClient() {
  return new google.auth.OAuth2(
    settings.oauth2.clientId,
    settings.oauth2.clientSecret,
    settings.oauth2.redirectUri,
  );
}

function getLoginUrl() {
  return newClient().generateAuthUrl({
    access_type: 'offline',
    prompt: 'select_account consent',
    scope: ['https://www.googleapis.com/auth/calendar.events'],
  });
}

exports.authenticate = function(request, response) {
  const client = newClient();
  client.getToken(request.query.code, (err, token) => {
    if (err || !token.refresh_token) {
      response.status(403).send(err);
      return;
    }

    const id = crypto.randomBytes(16).toString('base64');
    const key = datastore.key(['User', id]);
    const entity = { Token: token.refresh_token };
    datastore.set(key, entity).then(() => {
      request.session.id = id;
      response.redirect('/index.html');
    });
  });
};

exports.withCredentials = function(request, response, callback) {
  if (!request.session.id) {
    response.send({url: getLoginUrl()});
    return;
  }

  const key = datastore.key(['User', request.session.id]);
  datastore.get(key).then(entity => {
    if (!entity || !entity.Token) {
      response.send({url: getLoginUrl()});
      return;
    }

    const token = entity.Token;
    const client = newClient();
    client.setCredentials({refresh_token: token});
    callback(client);
  });
};

exports.logout = function(request, response) {
  request.session.id = null;
  response.send('You have been logged out');
};
