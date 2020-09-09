/**
 * All strings used in the UI.
 */
export interface StringAssets {
  submitButton: string;
  consentMessageTitle: string;
  consentMessageBodyTexts: string[];
  languageSelect: string;
  consentButton: string;
}

/**
 * The list of language code which supported by G Suite. Google meet supports
 * all of them.
 * See more information at
 * https://developers.google.com/admin-sdk/directory/v1/languages.
 * NOTE: right to left languages are not supported yet so they are excluded.
 */
export enum LanguageCode {
  Amharic = 'am',
  // Arabic='ar',
  Basque = 'eu',
  Bengali = 'bn',
  'English(UK)' = 'en-GB',
  'Portuguese(Brazil)' = 'pt-BR',
  Bulgarian = 'bg',
  Catalan = 'ca',
  Cherokee = 'chr',
  Croatian = 'hr',
  Czech = 'cs',
  Danish = 'da',
  Dutch = 'nl',
  'English(US)' = 'en',
  Estonian = 'et',
  Filipino = 'fil',
  Finnish = 'fi',
  French = 'fr',
  German = 'de',
  Greek = 'el',
  Gujarati = 'gu',
  // Hebrew = 'iw',
  Hindi = 'hi',
  Hungarian = 'hu',
  Icelandic = 'is',
  Indonesian = 'id',
  Italian = 'it',
  Japanese = 'ja',
  Kannada = 'kn',
  Korean = 'ko',
  Latvian = 'lv',
  Lithuanian = 'lt',
  Malay = 'ms',
  Malayalam = 'ml',
  Marathi = 'mr',
  Norwegian = 'no',
  Polish = 'pl',
  'Portuguese(Portugal)' = 'pt-PT',
  Romanian = 'ro',
  Russian = 'ru',
  Serbian = 'sr',
  'Chinese(PRC)' = 'zh-CN',
  Slovak = 'sk',
  Slovenian = 'sl',
  Spanish = 'es',
  Swahili = 'sw',
  Swedish = 'sv',
  Tamil = 'ta',
  Telugu = 'te',
  Thai = 'th',
  'Chinese(Taiwan)' = 'zh-TW',
  Turkish = 'tr',
  // Urdu='ur',
  Ukrainian = 'uk',
  Vietnamese = 'vi',
  Welsh = 'cy',
}

/**
 * Returns the language name specified by |code| in native writing.
 * Source: https://en.wikipedia.org/wiki/List_of_language_names
 */
export function getNativeName(code: LanguageCode): string {
  switch (code) {
    case LanguageCode.Amharic:
      return 'ኣማርኛ';
    case LanguageCode.Basque:
      return 'Euskara';
    case LanguageCode.Bengali:
      return 'বাংলা';
    case LanguageCode['English(UK)']:
      return 'English(UK)';
    case LanguageCode['Portuguese(Brazil)']:
      return 'Português';
    case LanguageCode.Bulgarian:
      return 'български език';
    case LanguageCode.Catalan:
      return 'Català';
    case LanguageCode.Cherokee:
      return 'ᏣᎳᎩ, ᏣᎳᎩ ᎦᏬᏂᎯᏍᏗ';
    case LanguageCode.Croatian:
      return 'Hrvatski';
    case LanguageCode.Czech:
      return 'Český Jazyk, Čeština';
    case LanguageCode.Danish:
      return 'Dansk';
    case LanguageCode.Dutch:
      return 'Nederlands';
    case LanguageCode['English(US)']:
      return 'English(US)';
    case LanguageCode.Estonian:
      return 'Eesti';
    case LanguageCode.Filipino:
      return 'Wikang Filipino';
    case LanguageCode.Finnish:
      return 'Suomi';
    case LanguageCode.French:
      return 'Français';
    case LanguageCode.German:
      return 'Deutsch';
    case LanguageCode.Greek:
      return 'Ελληνικά';
    case LanguageCode.Gujarati:
      return 'ગુજરાતી';
    case LanguageCode.Hindi:
      return 'हिन्दी';
    case LanguageCode.Hungarian:
      return 'Magyar';
    case LanguageCode.Icelandic:
      return 'Íslenska';
    case LanguageCode.Indonesian:
      return 'Bahasa Indonesia';
    case LanguageCode.Italian:
      return 'Italiano';
    case LanguageCode.Japanese:
      return '日本語';
    case LanguageCode.Kannada:
      return 'ಕನ್ನಡ';
    case LanguageCode.Korean:
      return '조선말, 한국어';
    case LanguageCode.Latvian:
      return 'latviešu';
    case LanguageCode.Lithuanian:
      return 'Lietuvių';
    case LanguageCode.Malay:
      return 'بهاس ملايو/Bahasa Melayu';
    case LanguageCode.Malayalam:
      return 'മലയാളം';
    case LanguageCode.Marathi:
      return 'मराठी';
    case LanguageCode.Norwegian:
      return 'Norsk';
    case LanguageCode.Polish:
      return 'Język polski';
    case LanguageCode['Portuguese(Portugal)']:
      return 'Português';
    case LanguageCode.Romanian:
      return 'Română';
    case LanguageCode.Russian:
      return 'Русский';
    case LanguageCode.Serbian:
      return 'Српски';
    case LanguageCode['Chinese(PRC)']:
      return '中文(简体)';
    case LanguageCode.Slovak:
      return 'Slovenčina';
    case LanguageCode.Slovenian:
      return 'Slovenščina';
    case LanguageCode.Spanish:
      return 'Español';
    case LanguageCode.Swahili:
      return 'Kiswahili';
    case LanguageCode.Swedish:
      return 'Svenska';
    case LanguageCode.Tamil:
      return 'தமிழ்';
    case LanguageCode.Telugu:
      return 'తెలుగు';
    case LanguageCode.Thai:
      return 'ภาษาไทย';
    case LanguageCode['Chinese(Taiwan)']:
      return '中文(繁体)';
    case LanguageCode.Turkish:
      return 'Türkçe';
    case LanguageCode.Ukrainian:
      return 'Українська';
    case LanguageCode.Vietnamese:
      return 'Tiếng Việt Nam';
    case LanguageCode.Welsh:
      return 'Cymraeg';
    default:
      checkExhaustive(code);
  }
}

function checkExhaustive(v: never): never {
  throw new Error('Unexpected Language Code');
}
