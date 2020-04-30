var languageAssets = {
  // English
  en: {
    continueButton: "Continue",
    welcomeMessage: "Welcome",
    waitingRoomMessage: "Your appointment will begin soon.<br />Thank you for your patience.",
    consentMessage: "<h1>Before you enter the waiting room, please review the following information about televisits, via MyChart:</h1>" +        
      "<ol>" +
        "<li>Please take a moment to see that you are in a private location where you cannot be unintentionally overheard. If you cannot be in a private location, you can decide at any time whether to continue this visit or to end it.</li>" +
        "<li>You understand there are potential risks to this technology, including interruptions, unauthorized access and technical difficulties.  You or your provider may need to discontinue the televisit at any time.</li>" +
        "<li>By continuing to participate in this telehealth visit, you are providing verbal consent for treatment.</li>" +
      "</ol>",
    languageSelect: "Select your language:"
  },
  // Spanish
  es: {
    continueButton: "Seguir",
    welcomeMessage: "Bienvenidos",
    waitingRoomMessage: "Tu cita comenzará pronto.<br/>agradecemos tu paciencia.",
    consentMessage: "<h1>Antes de entrar a la sala de espera, revise la siguiente información sobre las televisitas, a través de MyCHArt:</h1>" +        
      "<ol>" +
        "<li>Asegúrese de estar en un lugar privado, donde no lo puedan escuchar sin querer. Si no puede estar en un lugar privado, puede decidir continuar con la visita o terminarla en cualquier momento.</li>" +
        "<li>Usted entiende que existen posibles riesgos con esta tecnología, como interrupciones, acceso no autorizado y dificultades técnicas. Es posible que usted o su proveedor tengan que interrumpir la televisita en cualquier momento.</li>" +
        "<li>Al continuar con esta televisita, usted da su consentimiento verbal para el tratamiento.</li>" +
      "</ol>",
    languageSelect: "Elige tu idioma:"
  }
}

function getAssetsForLanguage(languageId) {
  var languageAsset = languageAssets[languageId];
  if (typeof languageAsset === 'undefined') { return languageAssets['en'] }
  return languageAsset;
}
