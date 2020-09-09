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

const net = require('net');
const settings = require('./settings.json'); 
const format = require('date-format');

const startBlock = '\x0b'
const endBlock   = '\x1c'
const cr = '\x0d'

exports.setAppointmentStatusArrived = (appointmentId, patientId, patientName) => {
    var msg =`MSH|^~\\&|Google|CHA|EPIC|CHA|${format('yyyyMMddhhmmss')}||SIU^S14|1058|P|2.6||${cr}`
           + `SCH||||||^Google Update||||||||||^Both parties arrived||||^Google|||||6^ARRIVED${cr}`
           + `PID|||${patientId}||${patientName}|||||||||||||||||||||||||${cr}`
           + `PV1|||||||||||||||||||${appointmentId}|||||||||||||||||${cr}`
    return mllpSend(msg);
}

const mllpSend = (msg) => {
    var socket = new net.Socket();
    return new Promise((res, rej) => {
        socket.connect(settings.mllpPort, settings.mllpIpAddress, function() {
            socket.write(`${startBlock}${msg}${endBlock}${cr}`);
        });

        socket.on('data', function(data) {
            // TODO: Validate response.
            socket.destroy();
            res();
        });

        socket.on('error', function(err) {
            rej(err);
        })
    })
}
