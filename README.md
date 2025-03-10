# MonkeyCode - A typing test CLI program but you get to code 🧑🏻‍💻 >_

MonkeyCode is a terminal-based typing test tool for developers inspired by the popular monkeytype typing app, where programmers practice typing by writing actual code in their favorite programming languages. Choose from languages like Go, C/C++, Python, Java, JavaScript, TypeScript etc. and challenge yourself with real coding tasks while tracking your Words Per Minute (WPM) and Accuracy.

Monkeycode is not just a platform for practicing typing; it's a new way to learn a new programming language or refresh your knowledge of an old one. By providing carefully curated code snippets, Monkeycode helps users reinforce syntax, improve coding speed, and build muscle memory in a hands-on, interactive manner.

The project also includes an SSH server directory, allowing users to connect and use the app remotely over SSH.

![Demo](demo.gif)

## Why? 
### Because, as a programmer, I have been wanting something like this for so long.

Typing speed test apps such as monkeytype, typingspeedtest, 10fastfingers, etc. are often limited to plain text, missing the unique challenges of coding syntax. I can type fast (avg. 100wpm) on these normal typing speed app, but that doesn't translate to wrting code since you don't get to practice typing special characters such as: (), {}, ;, <>, [], etc. and other unique combinations often only show in code, not just a text document and I don't plan to write word docs for a living (jk).

MonkeyCode fills this gap by providing a typing test designed specifically for developers, with realistic code snippets that help improve typing speed and accuracy in real-world coding scenarios. As a CLI tool, it can help you and me be more engaged with the terminal, something that I think is super important for a programmer.

On top of the program that you can run natively on your machine, with the included SSH server setup, MonkeyCode enables a collaborative and shared experience, making it ideal for coding bootcamps, hackathons, and remote practice.

## Installation (Local)
```bash
go install github.com/han-nwin/monkeycode/local@latest
mv "$(go env GOPATH)/bin/local" "$(go env GOPATH)/bin/monkeycode"
```

## Usage
### Create a user profile
```bash
monkeycode -u <user-name>
```
### Run the program
```bash
monkeycode
```
### Switch profile or create new profile
```bash
monkeycode -u <user-name>
```
### Show current user
```bash
monkeycode -c
```
### Display Leaderboard
```bash
monkeycode -l
```


