const { SlashCommandBuilder } = require('@discordjs/builders');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('prune')
		.setDescription('Prunes messages from the desired channel.')
			.addIntegerOption(option =>
				option.setName('amount')
				.setDescription('The amount of messages to delete')
				.setRequired(true)),
	async execute(interaction) {
		const numberToDelete = interaction.options.getInteger('amount');
		const MAX_DEPTH = 5;
		let usersMsgs = (await interaction.channel.messages.fetch({ limit: 100 })).filter(m => m.author.id === interaction.user.id)
		let toDelete = usersMsgs.clone()
		let lastMsg = usersMsgs.last().id;
		if(!lastMsg) return;
		let curr = 1
		while (usersMsgs.size == 100 && curr !== MAX_DEPTH) {
			usersMsgs = (await interaction.channel.messages.fetch({ limit: 100, before: lastMsg })).filter(m => m.author.id === interaction.user.id)
			lastMsg = usersMsgs.last().id
			curr++;
			toDelete = toDelete.concat(usersMsgs)
		}
		console.log(toDelete)
		// await interaction.channel.bulkDelete(msgs)
		await interaction.reply({ content: `Deleted ${0} messages`, ephemeral: true });
	},
};