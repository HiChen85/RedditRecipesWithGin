let user = {
    username: "hccc",
    password: 'yxxx_hccc'
}
let token = '';

async function login() {
    // 不能直接运行脚本, 需要在浏览器中运行 fetch.
    let response = await fetch('http://localhost:8000/signin', {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(user)
    });
    return await response.json();
}

let loginBtn = document.querySelector('.login');
loginBtn.addEventListener('click', () => {
    login().then(data => {
        token = data.token;
        console.log(data.token);
    });
})

let listRecipeBtn = document.querySelector('.listRecipe');
let tag = document.querySelector('input');
let searchTagBtn = document.querySelector('.searchTag');

function showRecipes(recipes) {
    let oldBox = document.querySelector('.recipes');
    document.body.removeChild(oldBox);
    let box = document.createElement('div')
    box.className = 'recipes';
    document.body.appendChild(box);
    for (let i = 0; i < recipes.length; i++) {
        let title = document.createElement('h3');
        title.textContent = recipes[i].name;
        box.appendChild(title);

        let tags = document.createElement('ul');
        let tagTitle = document.createElement('h4');
        tagTitle.textContent = "Tags"
        tags.appendChild(tagTitle);
        for (let j = 0; j < recipes[i].tags.length; j++) {
            let li = document.createElement('li');
            li.textContent = recipes[i].tags[j];
            tags.appendChild(li);
        }
        box.appendChild(tags);

        let ingredients = document.createElement('ul');
        let ingredientTitle = document.createElement('h4');
        ingredientTitle.textContent = "Ingredients";
        ingredients.appendChild(ingredientTitle);
        for (let j = 0; j < recipes[i].ingredients.length; j++) {
            let li = document.createElement('li');
            li.textContent = recipes[i].ingredients[j];
            ingredients.appendChild(li);
        }
        box.appendChild(ingredients);

        let instructions = document.createElement('ol');
        let instrTitle = document.createElement('h4');
        instrTitle.textContent = "Instructions";
        instructions.appendChild(instrTitle);
        for (let j = 0; j < recipes[i].instructions.length; j++) {
            let li = document.createElement("li");
            li.textContent = recipes[i].instructions[j];
            instructions.appendChild(li);
        }
        box.appendChild(instructions);
    }
}


listRecipeBtn.addEventListener('click', () => {
    let header = {
        'Authorization': token,
    }
    fetch('http://localhost:8000/recipes', {
        method: 'GET',
        mode: 'cors',
        headers: header,
    }).then(response => response.json()).then(recipes => showRecipes(recipes))
});

searchTagBtn.addEventListener('click', () => {
    let header = {
        'Authorization': token,
    }
    fetch(`http://localhost:8000/recipes/search?tag=${tag.value}`, {
        method: "GET",
        mode: "cors",
        headers: header,
    }).then(response => response.json()).then(recipes => {
        console.log(recipes);
        showRecipes(recipes);
    }).catch(e => {
        console.log(e.message);
    })
})


