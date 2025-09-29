// Minimal JavaScript for char counter and other utilities

// Character counter for tweet composer
document.addEventListener('DOMContentLoaded', function() {
  const textareas = document.querySelectorAll('textarea[maxlength]');
  
  textareas.forEach(textarea => {
    const maxLength = textarea.getAttribute('maxlength');
    
    textarea.addEventListener('input', function() {
      const length = this.value.length;
      const counter = document.getElementById('char-count');
      
      if (counter) {
        counter.textContent = length;
        
        // Update counter color based on length
        counter.classList.remove('warning', 'danger');
        if (length > maxLength * 0.9) {
          counter.classList.add('danger');
        } else if (length > maxLength * 0.8) {
          counter.classList.add('warning');
        }
      }
    });
  });
});

// Dark mode toggle (optional)
function toggleTheme() {
  const html = document.documentElement;
  const currentTheme = html.getAttribute('data-theme');
  const newTheme = currentTheme === 'light' ? 'dark' : 'light';
  html.setAttribute('data-theme', newTheme);
  localStorage.setItem('theme', newTheme);
}

// Load saved theme
const savedTheme = localStorage.getItem('theme') || 'light';
document.documentElement.setAttribute('data-theme', savedTheme);