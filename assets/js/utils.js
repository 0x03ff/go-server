// Updated checkCSRFToken function
export function checkCSRFToken() {
    // Helper function to get a cookie by name
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
        return null;
    }

    const csrfTokenInput = document.getElementById('csrf_token');
    const csrfToken = csrfTokenInput ? csrfTokenInput.value : '';
    
    // Check if the CSRF cookie exists
    const csrfCookie = getCookie('csrf_token');
    
    // If cookie is missing
    if (!csrfCookie) {
        console.log('CSRF cookie missing - refreshing page');
        setTimeout(() => {
            location.reload();
        }, 2000);
        return;
    }
    
    // If token input is missing or empty
    if (!csrfToken || csrfToken.trim() === '') {
        console.log('CSRF token input missing - refreshing page');
        setTimeout(() => {
            location.reload();
        }, 2000);
        return;
    }
    
    // If token length is suspiciously short
    if (csrfToken.length < 32) {
        console.log('CSRF token appears invalid - refreshing page');
        setTimeout(() => {
            location.reload();
        }, 2000);
        return;
    }
}
