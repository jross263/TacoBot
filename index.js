require('dotenv').config()
const path = require('path');
const fs = require('fs')
const { Client, Intents, Collection } = require('discord.js');
const Cleanup = require("./cleanup.js");
const admin = require("firebase-admin");
const serviceAccount = JSON.parse(process.env.FIRESTORE_CONFIG)

admin.initializeApp({
	credential: admin.credential.cert(serviceAccount)
});

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });
client.commands = new Collection();

const commandsPath = path.join(__dirname, 'commands');
const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));
for(const file of commandFiles){
	const filePath = path.join(commandsPath, file);
	const command = require(filePath);
	client.commands.set(command.data.name, command);
}

const eventsPath = path.join(__dirname, 'events');
const eventFiles = fs.readdirSync(eventsPath).filter(file => file.endsWith('.js'));
for (const file of eventFiles) {
	const filePath = path.join(eventsPath, file);
	const event = require(filePath);
	if (event.once) {
		client.once(event.name, (...args) => event.execute(...args));
	} else {
		client.on(event.name, (...args) => event.execute(...args));
	}
}

client.login(process.env.TOKEN);

Cleanup(() => {
  client.destroy()
})