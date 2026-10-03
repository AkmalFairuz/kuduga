class SupportTicketModel {
  int id;
  String categoryId;
  String category;
  int status;
  int createdAt;
  int updatedAt;
  List<SupportTicketMessageModel>? messages;

  SupportTicketModel(
      {required this.id,
      required this.categoryId,
      required this.category,
      required this.status,
      required this.createdAt,
      required this.updatedAt,
      this.messages});

  factory SupportTicketModel.fromJson(Map<String, dynamic> json) =>
      SupportTicketModel(
        id: json["id"],
        categoryId: json["categoryId"],
        category: json["category"],
        status: json["status"],
        createdAt: json["createdAt"],
        updatedAt: json["updatedAt"],
        messages: ((json["messages"] ?? []) as List<dynamic>)
            .map((e) => SupportTicketMessageModel.fromJson(e))
            .toList(),
      );

  String statusText() {
    switch (status) {
      case 0:
        return "Buka";
      default:
        return "Ditutup";
    }
  }
}

class SupportTicketMessageModel {
  int id;
  String author;
  int role;
  String message;
  List<String> attachments;
  int createdAt;

  SupportTicketMessageModel({
    required this.id,
    required this.author,
    required this.role,
    required this.message,
    required this.createdAt,
    required this.attachments,
  });

  factory SupportTicketMessageModel.fromJson(Map<String, dynamic> json) =>
      SupportTicketMessageModel(
        id: json["id"],
        author: json["author"],
        role: json["role"],
        message: json["message"],
        createdAt: json["createdAt"],
        attachments: (json["attachments"] as List<dynamic>)
            .map((e) => e.toString())
            .toList(),
      );
}

class SupportTicketCategoryModel {
  String id;
  String label;

  SupportTicketCategoryModel({required this.id, required this.label});

  factory SupportTicketCategoryModel.fromJson(Map<String, dynamic> json) =>
      SupportTicketCategoryModel(id: json["id"], label: json["label"]);
}
