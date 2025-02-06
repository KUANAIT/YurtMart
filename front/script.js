async function fetchItems() {
    try {
        const response = await fetch("http://localhost:8086/items");
        const items = await response.json();
        console.log("Before display:", items);
        displayItems(items);
    } catch (error) {
        console.error("Error loading products:", error);
    }
}

function displayItems(items) {
    const container = document.getElementById("items-container");
    container.innerHTML = "";

    const categories = {};

    items.forEach(item => {
        if (!categories[item.category]) {
            categories[item.category] = [];
        }
        categories[item.category].push(item);
    });

    for (let category in categories) {
        const categorySection = document.createElement("div");
        categorySection.className = "category-section";

        const categoryHeader = document.createElement("h3");
        categoryHeader.innerText = category;
        categorySection.appendChild(categoryHeader);

        const categoryItemsContainer = document.createElement("div");
        categoryItemsContainer.className = "category-items";

        categories[category].forEach(item => {
            const div = document.createElement("div");
            div.className = "item";
            div.innerHTML = `
                <h3>${item.name}</h3>
                <p>Price: ${item.price}₸</p>
                <p>In stock: ${item.stock} </p>
                <p>Description: ${item.description}</p>
            `;
            categoryItemsContainer.appendChild(div);
        });

        categorySection.appendChild(categoryItemsContainer);
        container.appendChild(categorySection);
    }
}

fetchItems();

setInterval(fetchItems, 2000);
