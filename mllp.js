const net = require('net');
const settings = require('./settings.json'); 
const format = require('date-format');

exports.setAppointmentStatusArrived = (appointmentId, patientId, patientName) => {
    var msg = `MSH|^~\&|Google|CHA|EPIC|CHA|${format('yyyyMMddhhmmss')}||SIU^S14|1058|P|2.6||
    SCH||||||^Google Update||||||||||^Patient Joined Video||||^Google|||||6^ARRIVED
    PID|||${patientId}||${patientName}|||||||||||||||||||||||||
    PV1|||||||||||||||||||${appointmentId}|||||||||||||||||`
    var socket = new net.Socket();
    return new Promise((res, rej) => {
        socket.connect(settings.mllpPort, settings.mllpIpAddress, function() {
            socket.write(msg);
        });

        socket.on('data', function(data) {
            // TODO: Validate response.
            client.destroy();
            res();
        });

        socket.on('error', function(err) {
            rej(err);
        })
    })
}
