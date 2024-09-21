function addSubmitListener(id_prefix) {
    const formElem         = document.getElementById(id_prefix + 'apiForm');
    const methodElem       = document.getElementById(id_prefix + 'method');
    const endpointElem     = document.getElementById(id_prefix + 'endpoint');
    const requestBodyElem  = document.getElementById(id_prefix + 'requestBody');
    const responseCodeElem = document.getElementById(id_prefix + 'responseCode');
    const responseDataElem = document.getElementById(id_prefix + 'responseData');

    formElem.addEventListener('submit', async (e) => {
        e.preventDefault();
        const method = methodElem.value;
        const endpoint = endpointElem.value;
        const requestBody = requestBodyElem.value;

        try {
            const response = await fetch(endpoint, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: requestBody,
            });

            let data = await response.text();
            try {
                data = JSON.parse(data);
                data = JSON.stringify(data, null, 2);
            } catch (error) {}

            // Display response code and data
            responseCodeElem.textContent = `HTTP ${response.status}`;
            responseCodeElem.className = `response-code response-${Math.floor(response.status / 100)}`;

            responseDataElem.textContent = data;
        } catch (error) {
            responseCodeElem.textContent = 'Client-side error occurred';
            responseDataElem.textContent = error.toString();
        }
    });
}