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
