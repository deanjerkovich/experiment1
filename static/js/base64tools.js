// Base64 Tools Module
class Base64ToolsManager {
    constructor() {
        this.encodeCount = 0;
        this.decodeCount = 0;
        this.initializeEventListeners();
    }

    // Initialize event listeners
    initializeEventListeners() {
        // Use event delegation on the base64Content container
        const base64Content = document.getElementById('base64Content');
        if (base64Content) {
            base64Content.addEventListener('click', (e) => {
                const target = e.target;
                
                // Handle tab switching
                if (target.classList.contains('tab')) {
                    const tabName = target.getAttribute('data-tab');
                    if (tabName) {
                        this.switchTab(tabName);
                    }
                }
                
                // Handle action buttons
                if (target.hasAttribute('data-action')) {
                    const action = target.getAttribute('data-action');
                    switch (action) {
                        case 'encode':
                            this.encodeText();
                            break;
                        case 'decode':
                            this.decodeText();
                            break;
                        case 'copy':
                            const copyTarget = target.getAttribute('data-target');
                            if (copyTarget) {
                                this.copyToClipboard(copyTarget);
                            }
                            break;
                    }
                }
            });
        }

        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.key === 'Enter') {
                if (e.target.id === 'encodeInput') {
                    this.encodeText();
                } else if (e.target.id === 'decodeInput') {
                    this.decodeText();
                }
            }
        });
    }

    // Switch between encode and decode tabs
    switchTab(tabName) {
        // Hide all tab contents
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        
        // Remove active class from all tabs
        document.querySelectorAll('.tab').forEach(tab => {
            tab.classList.remove('active');
        });
        
        // Show selected tab content
        const selectedTab = document.getElementById(tabName + '-tab');
        if (selectedTab) {
            selectedTab.classList.add('active');
        }
        
        // Add active class to selected tab
        const tabButton = document.querySelector(`.tab[data-tab="${tabName}"]`);
        if (tabButton) {
            tabButton.classList.add('active');
        }
        
        // Hide any existing results
        this.hideResults();
        this.hideMessage();
    }

    // Hide all result sections
    hideResults() {
        const results = ['encodeResult', 'decodeResult'];
        results.forEach(id => {
            const element = document.getElementById(id);
            if (element) element.classList.add('hidden');
        });
    }

    // Encode text to base64
    async encodeText() {
        const input = document.getElementById('encodeInput').value.trim();
        
        if (!input) {
            this.showMessage('Please enter some text to encode', 'error');
            return;
        }

        try {
            const response = await fetch('/api/base64/encode', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ text: input }),
            });

            const data = await response.json();

            if (data.success) {
                document.getElementById('encodedText').textContent = data.data.encoded;
                document.getElementById('encodeResult').classList.remove('hidden');
                this.encodeCount++;
                this.updateStats();
                this.showMessage(data.message, 'success');
            } else {
                if (response.status === 401) {
                    this.showMessage('Authentication required. Please login again.', 'error');
                    if (window.authManager) {
                        window.authManager.forceShowLogin();
                    }
                } else {
                    this.showMessage(data.message || 'Encoding failed', 'error');
                }
            }
        } catch (error) {
            console.error('Encoding error:', error);
            this.showMessage('An error occurred during encoding', 'error');
        }
    }

    // Decode text from base64
    async decodeText() {
        const input = document.getElementById('decodeInput').value.trim();
        
        if (!input) {
            this.showMessage('Please enter some base64 text to decode', 'error');
            return;
        }
        
        try {
            const response = await fetch('/api/base64/decode', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ text: input }),
            });

            const data = await response.json();

            if (data.success) {
                document.getElementById('decodedText').textContent = data.data.decoded;
                document.getElementById('decodeResult').classList.remove('hidden');
                this.decodeCount++;
                this.updateStats();
                this.showMessage(data.message, 'success');
            } else {
                if (response.status === 401) {
                    this.showMessage('Authentication required. Please login again.', 'error');
                    if (window.authManager) {
                        window.authManager.forceShowLogin();
                    }
                } else {
                    this.showMessage(data.message || 'Decoding failed', 'error');
                }
            }
                 } catch (error) {
             console.error('Decoding error:', error);
             this.showMessage('An error occurred during decoding', 'error');
         }
    }

    // Copy text to clipboard
    async copyToClipboard(elementId) {
        const text = document.getElementById(elementId).textContent;
        
        try {
            await navigator.clipboard.writeText(text);
            this.showMessage('Copied to clipboard!', 'success');
        } catch (err) {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = text;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            this.showMessage('Copied to clipboard!', 'success');
        }
    }

    // Update statistics display
    updateStats() {
        const encodeCountElement = document.getElementById('encodeCount');
        const decodeCountElement = document.getElementById('decodeCount');
        
        if (encodeCountElement) encodeCountElement.textContent = this.encodeCount;
        if (decodeCountElement) decodeCountElement.textContent = this.decodeCount;
    }

    // Reset statistics
    resetStats() {
        this.encodeCount = 0;
        this.decodeCount = 0;
        this.updateStats();
    }

    // Clear all inputs and results
    clearAll() {
        const inputs = ['encodeInput', 'decodeInput'];
        inputs.forEach(id => {
            const element = document.getElementById(id);
            if (element) element.value = '';
        });
        
        this.hideResults();
        this.hideMessage();
    }

    // Safe message display that doesn't depend on MessageManager
    showMessage(text, type = 'success') {
        try {
            if (window.MessageManager) {
                window.MessageManager.show(text, type);
            } else {
                // Fallback message display
                const messageElement = document.getElementById('message');
                if (messageElement) {
                    messageElement.textContent = text;
                    messageElement.className = `message ${type}`;
                    messageElement.classList.remove('hidden');
                    
                    // Auto-hide after 5 seconds
                    setTimeout(() => {
                        messageElement.classList.add('hidden');
                    }, 5000);
                }
            }
        } catch (error) {
            console.error('Error showing message:', error);
            // Last resort: console log
            console.log(`[${type.toUpperCase()}] ${text}`);
        }
    }

    // Safe message hiding
    hideMessage() {
        try {
            if (window.MessageManager) {
                window.MessageManager.hide();
            } else {
                const messageElement = document.getElementById('message');
                if (messageElement) {
                    messageElement.classList.add('hidden');
                }
            }
        } catch (error) {
            console.error('Error hiding message:', error);
        }
    }
}

// Export for use in other modules
window.Base64ToolsManager = Base64ToolsManager;
