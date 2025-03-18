package bot

const (
	startMsgRu = `
	👋 Добро пожаловать в бота генерации контекстных предложений! 🎉
	Я помогу вам учить новые слова, создавая примеры предложений на основе введённых вами слов. Просто отправьте мне слово, и я сгенерирую предложения, чтобы вы могли увидеть его в контексте.
	Перед началом используйте команду /preferences, чтобы настроить язык и уровень сложности — так я смогу подбирать для вас наиболее полезные предложения.
	Удачи в изучении! 📚✨`
	startMsgEn = `
	👋 Welcome to the Context Sentence Generator Bot! 🎉
	I help you learn new words by generating example sentences based on the words you provide. Just send me a word, and I'll create sentences to help you understand it in context.
	Before you start, use the /preferences command to set your language and difficulty level so I can generate the most useful sentences for you.
	Happy learning! 📚✨`
	helpMsgRu = `
	Вот список доступных команд:
	/preferences – Настроить язык и уровень сложности, чтобы получать примеры предложений, соответствующие вашему уровню.
	/help – Показать этот список команд с пояснениями.
	/premium – Открыть неограниченный доступ к генерации предложений с подпиской на премиум.
	Нужна дополнительная помощь? Просто напишите мне - @dafraer`
	helpMsgEn = `
	Here’s a list of commands you can use:
	/preferences – Set your language and difficulty level to get sentences tailored to your learning needs.
	/help – View this list of commands and their explanations.
	/premium – Unlock unlimited access to sentence generation with a premium subscription.
	Need more assistance? Just send me a message - @dafraer`
	chooseLangRu     = "🌍 Пожалуйста, выберите название языка, который вы изучаете!"
	chooseLangEn     = "🌍 Please select the language you are learning!"
	chooseLevelRu    = "Пожалуйста, выберите уровень языка для ваших предложений!"
	chooseLevelEn    = "Please choose the language level for your sentences!"
	preferencesSetRu = `
	Всё готово! ✅
	Теперь вы можете отправлять слова, для которых хотите сгенерировать предложения. Просто вводите их по одному, и я всё сделаю!`
	preferencesSetEn = `
	Everything is set! ✅
	Now you can send the words for which you’d like to generate sentences. Just type them in one by one, and I’ll do the rest!`
	unknownCommandRu = "Sorry, I don't know this command"
	unknownCommandEn = "Извините, я не знаю такой команды"
	responseMsgRu    = "Вот ваше предложение и перевод на русский:\n``` %s```\n``` %s```"
	responseMsgEn    = "Here is your sentence and english translation:\n``` %s```\n``` %s```"
	tooLongMsgEn     = "Sorry, your message is too long"
	tooLongMsgRu     = "Извините, ваше сообщение слишком длинное"
)
