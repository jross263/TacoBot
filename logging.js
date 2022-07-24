const { getFirestore, FieldValue } = require('firebase-admin/firestore');

const db = getFirestore();

module.exports = async function log(interaction) {
  db.collection('commandLog').add({
    commandId: interaction.commandId,
    commandName: interaction.commandName,
    userId: interaction.user.id,
    timestamp: FieldValue.serverTimestamp()
  });
}