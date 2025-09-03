// Main Application Module
class App {
    constructor() {
        this.modules = {};
        this.initializeModules();
        this.start();
    }

    // Initialize all application modules
    initializeModules() {
        try {
            // Initialize message manager first (other modules depend on it)
            if (typeof MessageManager !== 'undefined') {
                this.modules.messageManager = new MessageManager();
                window.messageManager = this.modules.messageManager;
            } else {
                console.warn('MessageManager not available, using fallback message handling');
            }
            
            // Initialize forms manager
            this.modules.formsManager = new FormsManager();
            window.formsManager = this.modules.formsManager;
            
            // Initialize authentication manager
            this.modules.authManager = new AuthManager();
            window.authManager = this.modules.authManager;
            
            // Initialize navigation manager
            this.modules.navigationManager = new NavigationManager();
            window.navigationManager = this.modules.navigationManager;
            
            // Initialize base64 tools manager
            this.modules.base64ToolsManager = new Base64ToolsManager();
            window.base64ToolsManager = this.modules.base64ToolsManager;
            
            console.log('All modules initialized successfully');
        } catch (error) {
            console.error('Error initializing modules:', error);
        }
    }

    // Start the application
    async start() {
        try {
            console.log('Starting application...');
            
            // Check for existing session
            await this.modules.authManager.checkExistingSession();
            
            // Set up global error handling
            this.setupErrorHandling();
            
            // Set up performance monitoring
            this.setupPerformanceMonitoring();
            
            console.log('Application started successfully');
        } catch (error) {
            console.error('Error starting application:', error);
            console.error('Failed to start application');
        }
    }

    // Set up global error handling
    setupErrorHandling() {
        // Handle unhandled promise rejections
        window.addEventListener('unhandledrejection', (event) => {
            console.error('Unhandled promise rejection:', event.reason);
            console.error('An unexpected error occurred');
        });

        // Handle JavaScript errors
        window.addEventListener('error', (event) => {
            console.error('JavaScript error:', event.error);
            console.error('A JavaScript error occurred');
        });
    }

    // Set up performance monitoring
    setupPerformanceMonitoring() {
        // Monitor page load performance
        window.addEventListener('load', () => {
            if ('performance' in window) {
                const perfData = performance.getEntriesByType('navigation')[0];
                console.log('Page load time:', perfData.loadEventEnd - perfData.loadEventStart, 'ms');
            }
        });
    }

    // Get module instance
    getModule(moduleName) {
        return this.modules[moduleName];
    }

    // Restart the application
    async restart() {
        console.log('Restarting application...');
        
        try {
            // Clean up existing modules
            this.cleanup();
            
            // Reinitialize modules
            this.initializeModules();
            
            // Restart
            await this.start();
            
            console.log('Application restarted successfully');
        } catch (error) {
            console.error('Error restarting application:', error);
        }
    }

    // Clean up resources
    cleanup() {
        // Clear any intervals or timeouts
        if (this.modules.authManager) {
            // Stop session monitoring
            this.modules.authManager.currentUser = null;
        }
        
        // Clear modules
        this.modules = {};
    }

    // Get application status
    getStatus() {
        return {
            modules: Object.keys(this.modules),
            authStatus: this.modules.authManager ? !!this.modules.authManager.currentUser : false,
            timestamp: new Date().toISOString()
        };
    }
}

// Initialize the application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM ready, initializing application...');
    window.app = new App();
});

// Export for use in other modules
window.App = App;
