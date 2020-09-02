const axios = require('axios');

exports.checkFhirAuthorization = async (fhirUrl, fhirToken, encounterId) => {
    try {
        const res =  await axios.get(
            `https://${fhirUrl}/Encounter/${encounterId}`,
            { headers: { Authorization: `Bearer ${fhirToken}` }});
        if (res != 200) {
            return new Error(`cannot find encounter ${encounterId}`)
        }
    } catch (err) {
        return err
    };
}
