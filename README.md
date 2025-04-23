<br />
<div align="center">

<h3 align="center">WordBuddy Telegram Bot</h3>

  <p align="center">
    WordBuddy is a Telegram bot designed to help language learners by generating contextual example sentences and audio pronunciations for user-provided words — making vocabulary practice more natural and engaging.
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#try-wordbuddy-for-yourself">Try WordBuddy for Yourself</a></li>
        <li><a href="#getting-started">Getting started</a></li>
      </ul>
    </li>
    <li><a href="#features">Features</a></li>
    <li><a href="#under-the-hood">Under the hood</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

WordBuddy is a Telegram bot built to support language learners by generating contextual example sentences and audio for any user-provided word. It was created with a specific goal in mind: to simplify the process of making high-quality Anki flashcards.

Anki is a powerful tool for memorization, especially when used with sentence-based learning instead of isolated vocabulary. Learning new words in context helps reinforce meaning, grammar, and usage — all at once. However, creating example sentences and generating pronunciation audio manually can be time-consuming and tedious. That's where WordBuddy steps in.

With WordBuddy, you simply send a word to the bot, and it returns:

    📜 A natural, contextual sentence using the word

    🔊 A clear audio pronunciation of the sentence

This makes it incredibly easy to copy everything directly into your Anki deck — no need to research examples or create audio files yourself.

Whether you’re learning a new language from scratch or expanding your vocabulary, WordBuddy helps you streamline your workflow and focus on what really matters: meaningful, effective learning.




<!-- GETTING STARTED -->
## Getting Started

### Try WordBuddy for Yourself
WordBuddy is available on Telegram! Start learning now by chatting with [@sentencegenbot](https://t.me/sentencegenbot).  
For a full list of features, check out the [Features](#features) section.

### Run Locally 
You can run WordBuddy locally by following these steps:

#### 1. Set up Google Application Default Credentials 
Use [this](https://cloud.google.com/docs/authentication/application-default-credentials) instruction to do that

#### 2. Get telegram bot token and API keys 
- Get telegram bot token from [@BotFather](https://t.me/BotFather))
- Get narakeet API key [here](https://www.narakeet.com/) (Only used for Georgian language)
- Get gemini API key [here](https://aistudio.google.com/apikey)

<b>Keep in mind that APIs might not be free</b>

#### 2. Start the Bot using Go
Run the following command to start WordBuddy:

```sh
go run cmd/main.go <tg-bot-token> <gemini-api-key> <narakeet-api-key>
```  

Now your bot should be up and running locally!


<!-- FEATURES -->
## Features

- **Multilingual Support**  
  Choose from a wide range of languages including:
    - 🌐 Major languages: **English, Russian, Spanish, French, German, Turkish, Greek, Japanese, Korean, Arabic, Italian**
    - 🗣️ Smaller languages: **Georgian, Tatar**

- **Customizable Difficulty Levels**  
  Generate sentences tailored to your learning level — from **A1 (beginner)** all the way to **C2 (advanced)**.

- **Bilingual UI**  
  The bot interface is available in both **English** and **Russian**, making it accessible for a wider audience.

> ⚠️ **Note:** This bot uses AI to generate content. While it performs well in major languages, it may occasionally produce mistakes or inaccuracies — especially in smaller or less-resourced languages.

## Under the Hood

Here's a breakdown of the tech powering the bot:

- **Core Logic**  
  Written in **Go**, using the [`github.com/go-telegram/bot`](https://github.com/go-telegram/bot) library for seamless Telegram integration.

- **Language Model**  
  Utilizes [**Gemini**](https://gemini.google.com/app?hl=en), a powerful LLM (Large Language Model), to generate grammatically and contextually accurate sentences across a variety of languages.

- **Database**  
  All user data and state are managed through [**Google Firestore**](https://firebase.google.com/docs/firestore), ensuring speed, scalability, and reliability.

- **Audio Generation**
    - For major languages (e.g. English, Spanish, Japanese, etc.), audio is generated using the [**Google Text-to-Speech API**](https://cloud.google.com/text-to-speech).
    - For **Georgian**, audio is generated via [**Narakeet**](https://www.narakeet.com/languages/georgian-text-to-speech/#trynow).
    - For **Tatar**, audio is sourced from the [**ISSAI**](https://issai.nu.edu.kz/ru/tatartts-rus/) website.

- **Deployment**  
  The entire app is deployed on [**Google Cloud Run**](https://cloud.google.com/run), enabling fast, serverless, and scalable performance.








<!-- CONTACT -->
## Contact Info

Kamil Nuriev- [telegram](https://t.me/dafraer) - kdnuriev@gmail.com