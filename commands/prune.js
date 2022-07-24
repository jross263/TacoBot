const { SlashCommandBuilder } = require('@discordjs/builders');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('prune')
		.setDescription('Prunes messages from the desired channel. Can go back a max of 100 messages.')
		.addIntegerOption(option =>
			option.setName('amount')
				.setDescription('The amount of messages to delete')
				.setRequired(true))
		.addUserOption(option =>
			option.setName('user')
				.setDescription("User's message to delete")
				.setRequired(false)),
	async execute(interaction) {
		try {
			await interaction.deferReply({ ephemeral: true });

			const numberToDelete = interaction.options.getInteger('amount');
			const userToDelete = interaction.options.getUser('user')?.id ?? interaction.user.id;
			const isGuildOwner = (await interaction.guild.fetchOwner()).id === interaction.user.id

			if (numberToDelete > 100) return await interaction.editReply({ content: "Your amount is bigger than 100", ephemeral: true });
			if (!isGuildOwner && userToDelete !== interaction.user.id) return await interaction.editReply({ content: "This option is only available to server owner.", ephemeral: true });

			let messages = (await interaction.channel.messages.fetch({ limit: 100 })).filter(m => m.author.id === userToDelete)

			if (messages.size > numberToDelete) {
				const stopVal = messages.size - numberToDelete
				for (let i = 0; i < stopVal; i++) {
					messages.delete(messages.lastKey())
				}
			}

			await interaction.channel.bulkDelete(messages)

			await interaction.editReply({ content: `Deleted ${messages.size} messages`, ephemeral: true });
		} catch (e) {
			await interaction.editReply({ content: `Error completing command, please try again.`, ephemeral: true });
		}
	},
};