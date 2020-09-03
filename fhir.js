const axios = require('axios');

exports.checkFhirAuthorization = async (fhirUrl, fhirToken, encounterId) => {
    const res =  await axios.get(
        `${fhirUrl}/Encounter/${encounterId}`,
        { headers: { Authorization: `Bearer ${fhirToken}` }});
    if (res.status != 200) {
        throw new Error(`cannot find encounter ${encounterId}`)
    }
}
