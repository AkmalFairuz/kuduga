class PulsaCategory {
  String prefix;
  int categoryId;

  PulsaCategory(this.prefix, this.categoryId);

  factory PulsaCategory.fromJson(Map<String, dynamic> json) =>
      PulsaCategory(json["prefix"], json["categoryId"]);
}
