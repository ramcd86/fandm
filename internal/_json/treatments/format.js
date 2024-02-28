const fs = require("fs");

fs.readFile("treatments_raw.json", "utf8", (err, data) => {
  if (err) {
    console.error(err);
    return;
  }

  // Process the data here
  const formattedData = {};
  let db_name = "";
  let db_desc = "";

  JSON.parse(data).forEach((treatment) => {
    if (treatment?.brandName) {
      db_name = treatment.brandName.toLowerCase().replace(/[\s+]/g, "_");
      db_desc =
        treatment.brandName.charAt(0).toUpperCase() +
        treatment.brandName.slice(1);
    } else {
      db_name = treatment.brandName.toLowerCase().replace(/[\s+]/g, "_");
    }
    db_desc =
      db_desc +
      " / " +
      treatment.genericName.charAt(0).toUpperCase() +
      treatment.genericName.slice(1);
    formattedData[db_name] = db_desc;
  });

  // write formattedData as JSON to treatments.json
  fs.writeFile("treatments.json", JSON.stringify(formattedData), (err) => {
    if (err) {
      console.error(err);
      return;
    }
  });
});
