const fs = require("fs");

fs.readFile("food_raw.json", "utf8", (err, foods) => {
  if (err) {
    console.error(err);
    return;
  }

  const formattedIngredients = {};

  JSON.parse(foods).forEach((food) => {
    food.ingredients.forEach((ingredient) => {
      const formattedIngredient = ingredient
        .toLowerCase()
        .replace(/[\s-]/g, "_")
        .replace(/['®%]/g, "");
      formattedIngredients[formattedIngredient] =
        ingredient.charAt(0).toUpperCase() + ingredient.slice(1);
    });
  });

  fs.writeFile("foods.json", JSON.stringify(formattedIngredients), (err) => {
    if (err) {
      console.error(err);
      return;
    }
  });
});
