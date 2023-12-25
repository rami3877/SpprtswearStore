// let cartIcon = document.querySelector('#cart-icon');
// let cart = document.querySelector('.cart');
// let closeCart = document.querySelector('#close-cart');

// // open cart
// cartIcon.onclick = () => {
//     cart.classList.add('active');
// }
// // close cart
// closeCart.onclick = () => {
//     cart.classList.remove('active');
// }

// let com = document.querySelector('#com');
// let comments = document.querySelector('.comments');
// let closeCom = document.querySelector('#close-com');
// // open comments
// //com.onclick = () => {
// //    comments.classList.add('active');
// //}
// // close comments
// // closeCom.onclick = () => {
// //     comments.classList.remove('active');
// // }

// var ProductChick = null
// var Orders = []

// class Product {
//     constructor(product) {
//         this.id = product.id
//         this.price = Number(product.querySelector("#price").innerText)
//         this.Description = product.querySelector("#Description").innerText
//         this.color = null
//         this.size = null
//         this.img = product.querySelector("#imgProduct").src
//         this.qty = 1
//         product.querySelectorAll(".color").forEach(element => {
//             if (element.style.scale == 1.3) {
//                 this.color = element.style.background
//             }
//         });
//         product.querySelectorAll("#size").forEach(element => {
//             if (element.style.background == "black") {
//                 this.size = element.innerText
//             }
//         });
//     }
// }
// function chosenProduct(product) {
//     ProductChick = new Product(product)
// }

// function showColors(sizeName) {

//     var listColor = sizeName.parentElement.parentElement.children[0].children;
//     var sizeList = sizeName.parentElement.children
//     for (const node of sizeList) {
//         if (node.innerText == sizeName.innerText) {
//             node.style.background = "black"
//         } else {
//             node.style.background = null
//         }

//     }
//     for (const node of listColor) {
//         if (node.id === sizeName.innerText) {
//             node.style.display = "block";
//         } else {
//             node.style.display = "none";
//             node.style.scale = 1
//         }
//     }

// }
// function checkColor(colorIn) {
//     var listColor = colorIn.parentElement.children
//     for (const colors of listColor) {
//         if (colorIn == colors) {
//             colors.style.scale = 1.3
//         } else {
//             colors.style.scale = 1
//         }
//     }
// }







// function inputChange(Value, defaultValue, input) {
//     var Price = Number(input.parentElement.children[1].children[0].children[1].innerText)
//     if ((defaultValue - Value) < 0) {
//         Price *= Value - defaultValue
//         document.getElementById("total-price").innerText = Number(document.getElementById("total-price").innerText) + Price
//     } else {
//         Price *= defaultValue - Value
//         document.getElementById("total-price").innerText = Number(document.getElementById("total-price").innerText) - Price
//     }
//     input.defaultValue = Value
// }


// function setCookie(name, value, days) {
//     const data = new Date();
//     data.setTime(data.getTime() + (days * 24 * 60 * 60 * 1000))
//     let expires = "expires=" + data.toUTCString();
//     document.cookie = `${name}=${value};${expires};path=${document.location.pathname}`
// }

// function deleteCookie(name) {
//     setCookie(name, null, null)
// }
// function getCoookie(name) {
//     const cookie = decodeURIComponent(document.cookie)
//     return cookie
// }

// /*
// <div class="cart-box">
//     <img src="https://www.nielsenanimal.com/wp-content/uploads/2018/09/mens-nike-athletic-dry-t-shirt-blue-t-shirts.jpg"alt="" class="cart-img">
//     <div class="detail-box">
//         <div class="cart-product-title">Product Shirt</div>
//         <div class="flex-details">
//             <span class="cart-price"><span>$</span><span>30</span></span>
//             <span id="size">XL</span>
//             <span class="colortest" id="color" style="background:red;"></span>
//         </div>
//         <input id="input" type="number" name="" value="1" min="1" class="cart-quantity" onchange="inputChange(this.value , this.defaultValue , this)" />
//     </div>
//     <i class='bx bx-trash' id="cart-remove" onclick="deleteItemFromCart(this) "></i>
// </div> 
// */

// /*
//  {
//     "idModel":4,
//     "color":"red", 
//     "size":"xl",
//     "Container":"newContainer",
//     "Kind":  "short"
// }             
// */
// function ItemCart(item) {

// }




// // function deleteItemFromCart(itemInCart) {
// //     let parentElement = itemInCart.parentElement
// //     let price = Number(parentElement.children[1].children[1].children[0].children[1].innerText) * Number(parentElement.children[1].children[2].value)

// //     document.getElementById("total-price").innerText = Number(document.getElementById("total-price").innerText) - price
// //     parentElement.remove()
// // }





//                 // let d = document.location.pathname.replace("/","").split("/")
//                 var ProductOrderaChick = null
//                 function ProductOrder(order) {
//                     ProductOrderaChick = Orders[order.id]
//                 }
//                 function deleteItemFromCart() {
//                     var element = document.getElementById("cart-contant")
//                     console.log(element)

//                 }
//                 function addToCart() {

//                     if (ProductChick == null || ProductChick.color == null || ProductChick.size == null) {
//                         return
//                     }
//                     let cartContant = document.getElementById("cart-contant").querySelectorAll(`#${ProductChick.id}`)
//                     if (cartContant.length != 0) {
//                         for (const product of cartContant) {
//                             console.log(product)
//                             if (product.querySelector("#color").style.background == ProductChick.color &&
//                                 product.querySelector("#size").innerText == ProductChick.size) {
//                                 product.querySelector("#input").value++
//                                 document.getElementById("total-price").textContent = Number(document.getElementById("total-price").textContent) + ProductChick.price

//                                 return
//                             }
//                         }
//                     }

//                     var addToCarItem = `
//                            <div class="cart-box" id="${ProductChick.id}" onclick="ProductOrder(this)">
//                                <img src="${ProductChick.img}" alt="" class="cart-img">
//                                <div class="detail-box">
//                                    <div class="cart-product-title">${ProductChick.Description}</div>
//                                    <div class="flex-details">
//                                        <span class="cart-price">
//                                            <span>$</span>
//                                            <span>${ProductChick.price}</span>
//                                        </span><span id="size">${ProductChick.size}</span>
//                                        <span class="colortest" id="color" style="background:${ProductChick.color};"></span>
//                                    </div><input id="input" type="number" name="" value="1" min="1"
//                                        onchange="inputChange(this.value,this.defaultValue,this)" class="cart-quantity" />
//                                </div><i class='bx bx-trash' id="cart-remove" onclick="deleteItemFromCart()"></i>
//                            </div>`
//                     document.getElementById("cart-contant").innerHTML += addToCarItem
//                     document.getElementById("total-price").textContent = Number(document.getElementById("total-price").textContent) + ProductChick.price
//                     Orders[Orders.length] = ProductChick
//                     ProductChick = null

//                 }


//                 /*
                
//                                 // this.imgHtml = document.createElement("img")
//                                 // this.imgHtml.setAttribute("src", this.img)
//                                 // this.imgHtml.setAttribute("class", "cart-img")
//                                 // let detail_box = document.createElement("div")
//                                 // detail_box.setAttribute("class" ,  "detail-box")
//                                 // this.DescriptionHtml = document.createElement("div")
//                                 // this.DescriptionHtml.setAttribute("class" , "cart-product-title")
//                                 // this.DescriptionHtml.textContent = this.Description
//                                 // let flex_details = document.createElement("div")
//                                 // flex_details.setAttribute("class" ,"flex-details")
//                                 // let cart_price = document.createElement("span")
//                                 // cart.setAttribute("class" , "cart-price")
//                                 // let span1 = document.createElement("span")
//                                 // let span2 = document. 
                
//                 */
//   //     return
//         //     let status = 0
//         //     fetch("/user/buy", {
//         //         method: "POST",
//         //         body: JSON.stringify(OrdersToServer),
//         //         headers: {
//         //             "Content-Type": "application/json",
//         //         },
//         //     }).then(re => {
//         //         status = re.status
//         //         return re.text()
//         //     }).then((data) => {
//         //         if (status == 301) {
//         //             window.location = data
//         //         } else {
                    
//         //         }
//         //     })

//             // static reloade(e) {
//             //     let newone = new Product(null)
//             //     newone.id = e.id
//             //     newone.price = e.price
//             //     newone.Description = e.Description
//             //     newone.color = e.color
//             //     newone.size = e.size
//             //     newone.img = e.img
//             //     newone.qty = e.qty
//             //     newone.elementHtml = document.createElement("div")
//             //     newone.elementHtml.setAttribute("class", "cart-box")
//             //     newone.elementHtml.innerHTML = ` <img src="${newone.img}"alt="" class="cart-img">
//             //                     <div class="detail-box">
//             //                         <div class="cart-product-title">${newone.Description}t</div>
//             //                         <div class="flex-details">
//             //                             <span class="cart-price"><span>$</span><span>${newone.price}</span></span>
//             //                             <span id="size">${newone.size}</span>
//             //                             <span class="colortest" id="color" style="background:${newone.color};"></span>
//             //                         </div>
//             //                         <input id="input" type="number" name="" value="${newone.qty}" min="1" class="cart-quantity" />
//             //                     </div>
//             //                     <i class='bx bx-trash' id="cart-remove"></i>
//             //                 `
//             //     newone.elementHtml.querySelector("#input").addEventListener("change", e => {
//             //         var total = Number(document.getElementById("total-price").textContent)
//             //         if ((newone.qty - Number(e.target.value)) < 0) {
//             //             total += (Number(e.target.value) - newone.qty) * newone.price
//             //         } else if ((newone.qty - Number(e.target.value)) > 0) {
//             //             total -= (newone.qty - Number(e.target.value)) * newone.price
//             //         }

//             //         document.getElementById("total-price").textContent = total
//             //         newone.qty = Number(e.target.value)
//             //     })
//             //     newone.elementHtml.querySelector("#cart-remove").addEventListener("click", e => {
//             //         var total = Number(document.getElementById("total-price").textContent) - (newone.qty * newone.price)
//             //         document.getElementById("total-price").textContent = total
//             //         newone.elementHtml.remove()
//             //         for (let i = 0; i < Order.length; i++) {
//             //             if (Order[i] == newone) {
//             //                 Order.splice(i, 1)
//             //             }
//             //         }
//             //     })

//             //     document.getElementById("total-price").textContent = Number(document.getElementById("total-price").textContent) + newone.price
//             //     document.getElementById("cart-contant").appendChild(newone.elementHtml)
//             //     Order[Order.length] = newone
//             // }
// // listOrderOld.forEach(e => {
//         //     Product.reloade(e)

//         // })
function Getproducts(){
    fetch('products.html')
    .then(response => {
        return response.text();
    })
    .then(html => {
        // Create a temporary container to parse the HTML content
        const tempContainer = document.createElement('div');
        tempContainer.innerHTML = html;

        // Extract product information from the container
        const products = Array.from(tempContainer.querySelectorAll('.product')).map(productElement => {
            return {
                id: productElement.getAttribute('data-id'),
                name: productElement.querySelector('h3').textContent,
                description: productElement.querySelector('p:nth-of-type(1)').textContent,
                price: productElement.querySelector('p:nth-of-type(2)').textContent.replace('Price: $', '')
            };
        });

        // Display the products
        displayProducts(products);
    })

// Function to display products in the HTML
function displayProducts(products) {
    const productListElement = document.getElementById('productList');

    // Create HTML to display each product
    const productHTML = products.map(product => {
        return `<div>
                    <h3>${product.name}</h3>
                    <p>${product.Description}</p>
                    <p>Price: $${product.price}</p>
                </div>`;
    }).join('');

    // Insert the product HTML into the productList element
    productListElement.innerHTML = productHTML;
}
}
function postbuy(){
    // Define the product data
const productData = {
    "idModel": 4,
    "color": "red",
    "size": "xl",
    "Container": "newContainer",
    "Kind": "short"
};

// Convert the product data into an array of objects
const productArray = Object.entries(productData).map(([key, value]) => ({ id: key, value }));

// Use the fetch function to make a POST request
fetch("/api/buy", {
    method: "POST",
    headers: {
        "Content-Type": "application/json",
    },
    body: JSON.stringify(productArray),
})
    .then(response => {return response.json();
    }).then(data => {
        // Handle the response data here, if needed
        console.log(data);
    })
   
}