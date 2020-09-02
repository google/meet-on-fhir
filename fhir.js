var fhir = require('fhir.js');

var client = mkFhir({
    baseUrl: 'http://try-fhirplace.hospital-systems.com'
});

client
    .search( {type: 'Patient', query: { 'birthdate': '1974' }})
    .then(function(res){
        var bundle = res.data;
        var count = (bundle.entry && bundle.entry.length) || 0;
        console.log("# Patients born in 1974: ", count);
    })
    .catch(function(res){
        //Error responses
        if (res.status){
            console.log('Error', res.status);
        }

        //Errors
        if (res.message){
            console.log('Error', res.message);
        }
    });