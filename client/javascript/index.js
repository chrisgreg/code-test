function initialise() {
    const initialState = {
        startTime: null,
        sessionID: generateSessionId(),
        initialDimensions: getCurrentDimensions(),
    }

    return initialState;
}

function applyEventListeners(elements) {
    for (let elementIndex = 0; elementIndex < elements.length; elementIndex++) {
        const element = elements[elementIndex];
        
        element.addEventListener('keypress', (e) => {
            setInitialTime();
        })
        
        element.addEventListener('paste', (e) => {
            setInitialTime();
            const inputId = e.target.id;
            logToServer({
                eventType: "copyAndPaste",
                pasted: true,
                inputId,
            })
        })
    }

    var resizeId;
    window.addEventListener('resize', (resize) => {
        clearTimeout(resizeId);
        resizeId = setTimeout(() => {
            const newDimensions = getCurrentDimensions();
            window.state = {
                ...window.state,
                newDimensions,
            };
            logToServer({
                eventType: "resize",
                resizeFrom: window.state.initialDimensions,
                resizeTo: newDimensions,
            });
        }, 500);
    })

    const form = document.querySelector('.form-details');
    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const elapsed = Math.ceil(Date.now() - window.state.startTime)
        logToServer({
            eventType: "elapsedTime",
            time: elapsed
        })
    })
}

function generateSessionId() {
    return Math.random().toString(36).substring(7);
}

function setInitialTime() {
    if (window.state.startTime === null) {
        window.state.startTime = Date.now();
    }
}

function getCurrentDimensions() {
    return {
        width: window.outerWidth,
        height: window.outerHeight,
    }
}

function logToServer(data) {
    const dataToPost = {
        sessionID: window.state.sessionID,
        websiteUrl: window.location.href,
        ...data
    };

    return fetch('http://localhost:3030/log', {
        method: "post",
        body: JSON.stringify(dataToPost),
    });
}

document.addEventListener('DOMContentLoaded', () => {
    window.state = initialise();
    
    // Add Event listeners to each input
    const inputs = document.querySelectorAll('input');
    applyEventListeners(inputs);
})
