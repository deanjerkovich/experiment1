// Message Manager Module
class MessageManager {
    constructor() {
        this.messageElement = null;
        this.messageTimeout = null;
        this.initialize();
    }

    // Initialize the message system
    initialize() {
        this.messageElement = document.getElementById('message');
        if (!this.messageElement) {
            console.error('Message element not found');
        }
    }

    // Show a message
    static show(text, type = 'success', duration = 5000) {
        const messageElement = document.getElementById('message');
        if (!messageElement) {
            console.error('Message element not found');
            return;
        }

        // Clear any existing timeout
        if (MessageManager.instance && MessageManager.instance.messageTimeout) {
            clearTimeout(MessageManager.instance.messageTimeout);
        }

        // Set message content and type
        messageElement.textContent = text;
        messageElement.className = `message ${type}`;
        messageElement.classList.remove('hidden');

        // Auto-hide after duration (unless it's an error)
        if (type !== 'error' && duration > 0) {
            MessageManager.instance.messageTimeout = setTimeout(() => {
                MessageManager.hide();
            }, duration);
        }
    }

    // Hide the message
    static hide() {
        const messageElement = document.getElementById('message');
        if (messageElement) {
            messageElement.classList.add('hidden');
        }
    }

    // Show success message
    static success(text, duration = 5000) {
        MessageManager.show(text, 'success', duration);
    }

    // Show error message
    static error(text, duration = 0) {
        MessageManager.show(text, 'error', duration);
    }

    // Show info message
    static info(text, duration = 3000) {
        MessageManager.show(text, 'info', duration);
    }

    // Show warning message
    static warning(text, duration = 4000) {
        MessageManager.show(text, 'warning', duration);
    }

    // Clear message timeout
    static clearTimeout() {
        if (MessageManager.instance && MessageManager.instance.messageTimeout) {
            clearTimeout(MessageManager.instance.messageTimeout);
            MessageManager.instance.messageTimeout = null;
        }
    }
}

// Create singleton instance
MessageManager.instance = new MessageManager();

// Export for use in other modules
window.MessageManager = MessageManager;
