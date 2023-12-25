function PostVisa() {
    
    const Visanumber = document.getElementById("Visanumber").value
    console.log(Visanumber)
    const cvv = document.getElementById("cvv").value
    console.log(cvv)
    const expirydate = document.getElementById("expirydate").value
    console.log(expirydate)
    fetch("/user/visa", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"Visanumber":"${Visanumber}","cvv":"${cvv}","expirydate":"${expirydate}"}`,
    }).then(Response => Response.text()).then(data => {

    })
    const mobilenumber = document.getElementById("mobilenumber").value = Number
    console.log(mobilenumber)
    fetch("/user/phone", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `mobilenumber=${mobilenumber}`,
    }).then(Response => Response.text()).then(data => {

    })
    const cardholdername = document.getElementById("cardholdername").value
    console.log(cardholdername)
    fetch("/user/name", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `cardholdername=${cardholdername}`,
    }).then(Response => Response.text()).then(data => {

    })
    // Get values from the input fields
const city = document.getElementById('city').value;
const street = document.getElementById('street').value;
const apartment = document.getElementById('apartment').value;
const building = document.getElementById('building').value;

// Construct the request body
const requestBody = `city=${encodeURIComponent(city)}&street=${encodeURIComponent(street)}&apartment=${encodeURIComponent(apartment)}&building=${encodeURIComponent(building)}`;

// Use the fetch function to make a POST request
fetch("/user/address", {
    method: "POST",
    headers: {
        "Content-Type": "application/x-www-form-urlencoded",
    },
    body: requestBody,
})
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.text();
    })
    .then(data => {
        // Handle the response from the server
        console.log(data);
    })
    .catch(error => {
        // Handle errors that occurred during the fetch
        console.error('Fetch error:', error);
    }); 

}