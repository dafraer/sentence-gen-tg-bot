package text

type Messages struct {
	Start              map[string]string //Sent on /start command
	Help               map[string]string //Sent on /help command
	Lang               map[string]string //Sent when prompting user to choose the language
	Level              map[string]string //Sent when prompting user to choose language level (e.g. A1)
	PreferencesSet     map[string]string //Sent after user finishes set up
	UnknownCommand     map[string]string //Sent when receiving unknown command
	ResponseMsg        map[string]string //Sent when sending generated sentences to the user
	TooLong            map[string]string //Sent when message exceeds maxMessageLen set in bot.go
	BadRequest         map[string]string //Sent when unable to make sentences due to word being inappropriate or not existing
	Premium            map[string]string //Sent when user uses /premium command if they don't have premium yet
	LimitReached       map[string]string //Sent when user reaches free limit of 50 sentences per day
	PremiumTitle       map[string]string //Title of the message with the invoice and text of premium inline
	SuccessfulPayment  map[string]string //Sent when payment is successful
	FailedPayment      map[string]string //Sent when payment has failed
	PreferencesNotSet  map[string]string //Sent when user tries to generate sentences without setting the preferences
	AlreadyPremium     map[string]string //Sent when premium user tries to buy premium
	PremiumDescription map[string]string //Sent in the description of the invoice
}

// Load returns a Message object with all the message in russian and english
func Load() *Messages {
	var msgs Messages
	msgs.Start = map[string]string{
		"ru": `
👋 Добро пожаловать в бота генерации контекстных предложений! 🎉
Я помогу вам учить новые слова, создавая примеры предложений на основе введённых вами слов. Просто отправьте мне слово, и я сгенерирую предложения, чтобы вы могли увидеть его в контексте.
Перед началом используйте команду /preferences, чтобы настроить язык и уровень сложности — так я смогу подбирать для вас наиболее полезные предложения.
Удачи в изучении! 📚✨`,
		"en": `
👋 Welcome to the Context Sentence Generator Bot! 🎉
I help you learn new words by generating example sentences based on the words you provide. Just send me a word, and I'll create sentences to help you understand it in context.
Before you start, use the /preferences command to set your language and difficulty level so I can generate the most useful sentences for you.
Happy learning! 📚✨`,
	}
	msgs.Help = map[string]string{
		"ru": `
📌 Доступные команды:
✅ /preferences – Выберите язык и уровень сложности для персонализированных предложений.  
✅ /help – Посмотреть список команд и их описание.  
✅ /premium – Получите неограниченную генерацию предложений и доступ к более продвинутой модели ИИ для еще лучших результатов.  
Нужна помощь? Напишите мне – @dafraer`,
		"en": `
📌 Available Commands:
✅ /preferences – Set your language and difficulty level for personalized sentences.  
✅ /help – View this list of commands and their explanations.  
✅ /premium – Get unlimited sentence generation and access a more advanced AI model for even better sentences.
Need help? Just send me a message – @dafraer`,
	}
	msgs.Lang = map[string]string{
		"ru": "🌍 Пожалуйста, выберите название языка, который вы изучаете!",
		"en": "🌍 Please select the language you are learning!",
	}
	msgs.Level = map[string]string{
		"ru": "Пожалуйста, выберите уровень языка для ваших предложений!",
		"en": "Please choose the language level for your sentences!",
	}
	msgs.PreferencesSet = map[string]string{
		"ru": `
Всё готово! ✅
Теперь вы можете отправлять слова, для которых хотите сгенерировать предложения. Просто вводите их по одному, и я всё сделаю!`,
		"en": `
Everything is set! ✅
Now you can send the words for which you’d like to generate sentences. Just type them in one by one, and I’ll do the rest!`,
	}
	msgs.UnknownCommand = map[string]string{
		"ru": "Sorry, I don't know this command",
		"en": "Извините, я не знаю такой команды",
	}
	msgs.ResponseMsg = map[string]string{
		//Response messages need escaping \ because they are parsed using telegram's Mark Down Parse mode
		"ru": "⚠️ Обратите внимание, что ИИ может иногда допускать ошибки\\.\nВот ваше предложение и перевод на русский:\n``` %s```\n``` %s```",
		"en": "⚠️ Please note that AI may occasionally make mistakes\\.\nHere is your sentence and english translation:\n``` %s```\n``` %s```",
	}
	msgs.TooLong = map[string]string{
		"ru": "Sorry, your word is too long",
		"en": "Извините, ваше слово слишком длинное",
	}
	msgs.BadRequest = map[string]string{
		"ru": "❌ Извините, я не могу составить предложение с этим словом. Оно может быть неуместным или отсутствовать в выбранном языке.",
		"en": "❌ Sorry, I can’t generate a sentence with that word. It may be inappropriate or not exist in the selected language.",
	}
	msgs.Premium = map[string]string{
		"ru": `
💎 Оформите Premium на 30 дней с неограниченным доступом! 💎
Генерируйте неограниченное количество предложений и используйте более продвинутую модель ИИ! 🚀 Premium поддерживает разработчика и покрывает расходы на API. 💙
Оформите подписку сейчас и улучшите процесс обучения! ✨`,
		"en": `
💎 Go Premium for 30 Days of Unlimited Access! 💎
Generate unlimited sentences, and enjoy a more advanced AI model! 🚀 Premium supports the creator and covers API costs. 💙
Upgrade now and enhance your learning experience! ✨`,
	}
	msgs.LimitReached = map[string]string{
		"ru": `
🚨 Дневной лимит исчерпан! 🚨
Вы использовали 50 бесплатных предложений. Хотите безлимитный доступ? 
Оформите Premium, чтобы продолжать обучение и поддержать бота! 💙`,
		"en": `
🚨 Daily Limit Reached! 🚨
You've used all 50 free sentences for today. Want unlimited access? 
Upgrade to Premium to keep learning and support the bot! 💙`,
	}
	msgs.PremiumTitle = map[string]string{
		"ru": "Подписка Premium - 30 дней",
		"en": "Premium Subscription - 30 days",
	}
	msgs.SuccessfulPayment = map[string]string{
		"ru": `
✅ Оплата успешно обработана! ✅
Теперь у вас неограниченный доступ к боту на 30 дней. Спасибо за поддержку! Желаю вам успехов в изучении языков! 📚✨`,
		"en": `
✅ Payment successfully processed! ✅
You now have unlimited access for 30 days. Thank you for your support! Wishing you success in your language learning journey! 📚✨`,
	}
	msgs.FailedPayment = map[string]string{
		"ru": "Извините, что-то пошло не так.😔 Напишите @dafraer для решения проблемы",
		"en": "Sorry, something went wrong.😔 Write @dafraer to solve your issue",
	}
	msgs.PreferencesNotSet = map[string]string{
		"ru": "⚙️Сначала настройте бота используя команду /preferences! Без этого бот не будет работать.",
		"en": "⚙️Set your preferences using /preferences command first! The bot won’t work until you do.",
	}
	msgs.AlreadyPremium = map[string]string{
		"en": `
🎉 Вы уже Premium пользователь! 🎉  
У вас осталось %d дней доступа к Premium. Спасибо за поддержку! 💙  
Наслаждайтесь неограниченной генерацией предложений!`,
		"ru": `
🎉 You're already a Premium user! 🎉
You currently have %d days of Premium access left. Thank you for supporting us! 💙  
Enjoy your unlimited sentence generation!`,
	}
	msgs.PremiumDescription = map[string]string{
		"ru": "Откройте неограниченную генерацию предложений и доступ к улучшенной ИИ модели.",
		"en": "Unlock unlimited sentence generation and access a better AI model",
	}
	return &msgs
}
