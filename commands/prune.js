const { SlashCommandBuilder } = require('@discordjs/builders');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('prune')
		.setDescription('Prunes messages from the desired channel. Can go back a max of 100 messages.')
			.addIntegerOption(option =>
				option.setName('amount')
				.setDescription('The amount of messages to delete')
				.setRequired(true)),
	async execute(interaction) {
		try{
			await interaction.deferReply({ ephemeral: true });
			const numberToDelete = interaction.options.getInteger('amount');
			if(numberToDelete > 100) interaction.editReply({ content: "Your amount is bigger than 100", ephemeral: true});

			let messages = (await interaction.channel.messages.fetch({ limit: 100 })).filter(m => m.author.id === interaction.user.id)

			if(messages.size > numberToDelete){
				const stopVal = messages.size - numberToDelete
				for(let i = 0; i < stopVal; i++){
					messages.delete(messages.lastKey())
				}
			}

			await interaction.channel.bulkDelete(messages)
			
			await interaction.editReply({ content: `Deleted ${messages.size} messages`, ephemeral: true });
		}catch(e){
			await interaction.editReply({ content: `Error completing command, please try again.`, ephemeral: true });
		}
	},
};