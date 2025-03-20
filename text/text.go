package text

type Messages struct {
	Start             map[string]string
	Help              map[string]string
	Lang              map[string]string
	Level             map[string]string
	PreferencesSet    map[string]string
	UnknownCommand    map[string]string
	ResponseMsg       map[string]string
	TooLong           map[string]string
	BadRequest        map[string]string
	Premium           map[string]string
	LimitReached      map[string]string
	PremiumTitle      map[string]string
	SuccessfulPayment map[string]string
	FailedPayment     map[string]string
}

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
		Вот список доступных команд:
		/preferences – Настроить язык и уровень сложности, чтобы получать примеры предложений, соответствующие вашему уровню.
		/help – Показать этот список команд с пояснениями.
		/premium – Открыть неограниченный доступ к генерации предложений с подпиской на премиум.
		Нужна дополнительная помощь? Просто напишите мне - @dafraer`,
		"en": `
		Here’s a list of commands you can use:
		/preferences – Set your language and difficulty level to get sentences tailored to your learning needs.
		/help – View this list of commands and their explanations.
		/premium – Unlock unlimited access to sentence generation with a premium subscription.
		Need more assistance? Just send me a message - @dafraer`,
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
		"ru": "Вот ваше предложение и перевод на русский:\n``` %s```\n``` %s```",
		"en": "Here is your sentence and english translation:\n``` %s```\n``` %s```",
	}
	msgs.TooLong = map[string]string{
		"ru": "Sorry, your message is too long",
		"en": "Извините, ваше сообщение слишком длинное",
	}
	msgs.BadRequest = map[string]string{
		"ru": "Извините, я не могу составить предложение с этим словом",
		"en": "Sorry, I can't make a sentence with that word.",
	}
	msgs.Premium = map[string]string{
		"ru": `
		💎 Оформите Premium и получите месяц безлимитного доступа! 💎
		Открывайте все функции и генерируйте бесконечные предложения! 🚀 Подписка покрывает расходы на API и поддерживает разработчика. 💙
		Оформите Premium и учитесь без ограничений! ✨`,
		"en": `
		💎 Go Premium for 1 Month of Unlimited Access! 💎
		Unlock all features and generate endless sentences! 🚀 Premium helps cover API costs and supports the creator. 💙
		Upgrade now and keep learning without limits! ✨`,
	}
	msgs.LimitReached = map[string]string{
		"ru": `
		🚨 Дневной лимит исчерпан! 🚨
		Вы использовали 50 бесплатных предложений. Хотите безлимитный доступ? 🚀
		Оформите Premium, чтобы продолжать обучение и поддержать бота! 💙✨`,
		"en": `
		🚨 Daily Limit Reached! 🚨
		You've used all 50 free sentences for today. Want unlimited access? 🚀
		Upgrade to Premium to keep learning and support the bot! 💙✨`,
	}
	msgs.PremiumTitle = map[string]string{
		"ru": "Премиум",
		"en": "Premium",
	}
	msgs.SuccessfulPayment = map[string]string{
		"ru": `
		✅ Оплата успешно обработана! ✅
		Теперь у вас неограниченный доступ к боту на 1 месяц. Спасибо за поддержку! Желаем вам успехов в изучении языков! 📚✨`,
		"en": `
		✅ Payment successfully processed! ✅
		You now have unlimited access for 1 month. Thank you for your support! Wishing you success in your language learning journey! 📚✨`,
	}
	msgs.FailedPayment = map[string]string{
		"ru": "Извините, что-то пошло не так. Напишите @dafraer для решения проблемы",
		"en": "Sorry, something went wrong. Write @dafraer to solve your issue",
	}
	return &msgs
}
