const lostIdentityLink = document.getElementById('ludo-lost-identity');
const popin = document.getElementById('ludo-popin');
const closePopinButton = document.getElementById('ludo-close-popin');
const overlay = document.getElementById('ludo-overlay');

// Function to open the popin
function openPopin() {
  popin.style.display = 'block';
  overlay.style.display = 'block';
}

// Function to close the popin
function closePopin() {
  popin.style.display = 'none';
  overlay.style.display = 'none';
}

// Add event listeners
lostIdentityLink.addEventListener('click', openPopin);
closePopinButton.addEventListener('click', closePopin);
overlay.addEventListener('click', closePopin);