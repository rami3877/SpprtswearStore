function Getorders(){
    fetch("/admin/orders", {
        method: "GET",
    }).then(Response => Response.json()).then(data => {
    }) 
}
// link products.js with manage.html ?? whete <script src="js/products.js"></script>

// Needs to be linked to a file manage.html file


function PostAdminProducts(){
  const containerValue = document.getElementById('container').value;
  const kindValue = document.getElementById('kind').value;
  console.log(containerValue)
  console.log(kindValue)
  fetch("/admin/product/container", {
    method: "POST",
    headers: {
        "Content-Type": "application/x-www-form-urlencoded",
    },
    body: `Container=${containerValue}`,
}).then(Response => Response.text()).then(data => {
  //call kind
    console.log(data)
    if (data == "Created"){
      fetch('/admin/product/kind', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: `{"kind":"${kindValue}", "container":"${containerValue}"}`,
      })
      .then(Response => Response.text()).then(data => {
        
          console.log(data)
      
      })}
})
}
  // POST MODEL 
// function postModel(){
//   const containerValue = document.getElementById('container').value;
//   const kindValue = document.getElementById('kind').value;
//   const discount = document.getElementById('discount').value;
//   const productname = document.getElementById('product-name').value;
//   const price = document.getElementById('price').value;
//   const image = document.getElementById('image').value;
//   const sizes = document.getElementById('size').value;
  
//   // Construct the postData object using the obtained values
//   const postData = {
//     containerValue: containerValue,
//     kind: kindValue,
//     model: {
//       sizes:  { xl: { red: 1 }, sizes },
//       price: price,
//       productname: productname, // Assuming productname is the description
//       discount: discount,
//       linkesImage: [image]
//     }
//   };

//   fetch("/admin/product/model", {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     body: JSON.stringify(postData),
//   })
//   .then(response => {return response.json();
    
//   })
//   console.log(postData)
// }

    /*const productData = {
        "idModel": 4,
        "color": "red",
        "size": "xl",
        "containerValue": "newContainer",
        "Kind": "short"
    };
    // Convert the product data into an array of objects
    const productArray = Object.entries(productData).map(([key, value]) => ({ id: key, value }));
    fetch("/admin/product/model", {
      
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(productArray),
    })
        .then(response => {return response.json();
    
        })*/


//const containerValue = document.getElementById('container').value;
  //const kindValue = document.getElementById('kind').value;
// GET CONTAINER
// function getcontainer(){
//     fetch("/admin/product/container", {
//         method: "GET"
//     }).then(re => re.ok).then(d => {
//             window.location.reload()
//     })
// }

//GET MODEL
function getmodel(){
    const deleteData = {
        container: "newContainer",
        kind: "short",
        id: 2
      };
      // Construct the URL with parameters
      fetch("/admin/product/model", {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(deleteData),
      })
      .then(response => response.json())
      .then(result => {
        // Handle the response result
        console.log(result);
      })
       
}

// function deleteKind(){
//   // delete kind 
//   fetch("/admin/product/kind", {
//     method: 'DELETE',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     body: `kind=${kindValue}`,
//   }).then(Response => Response.text()).then(data => {

//     console.log(data)

// })

// }

// function deletecontiners(){
//   const containerValue = document.getElementById('container').value;
// // DELETE container
// fetch("/admin/product/container", {
//     method: 'DELETE',
//     headers: {
//       "Content-Type": "application/json",
//   },
//   body: `Container=${containerValue}`,
//   }).then(Response => Response.text()).then(data => {

//     console.log(data)

// })
// }
// function Container(){
//   const containerValue = document.getElementById('container').value;
//   console.log(containerValue)
//   fetch("/admin/product/container", {
//     method: "POST",
//     headers: {
//         "Content-Type": "application/x-www-form-urlencoded",
//     },
//     body: `Container=${containerValue}`,
// }).then(Response => Response.text()).then(data => {

//     console.log(data)

// })
// }

// function Kind() {
//   const kindValue = document.getElementById('kind').value;

//   fetch('/admin/product/kind', {
//     method: 'POST',
//     headers: {
//       'Content-Type': 'application/json',
//     },
//     body: `kind=${kindValue}`,
//   }).then(Response => Response.json()).then(data => {
//       //console.log(data)
  
//   })
// }