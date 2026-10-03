import 'package:finance_app/model/support.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/utils/json.dart';
import 'package:finance_app/widgets/input/file_input.dart';
import 'package:http/http.dart';

class SupportService {
  static Future<List<SupportTicketCategoryModel>> getTicketCategories() async {
    var resp = await Server.get("/support/ticket/categories");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => SupportTicketCategoryModel.fromJson(e))
        .toList();
  }

  static Future<List<SupportTicketModel>> getTickets() async {
    var resp = await Server.get("/support/ticket/");
    return (jsonDec(resp) as List<dynamic>)
        .map((e) => SupportTicketModel.fromJson(e))
        .toList();
  }

  static Future<SupportTicketModel> getDetailedTicket(int ticketId,
      {int? messageAfterId}) async {
    var q = {"id": ticketId};
    if (messageAfterId != null) {
      q["messageAfterId"] = messageAfterId;
    }
    var resp = await Server.get("/support/ticket/detail", query: q);
    return SupportTicketModel.fromJson(jsonDec(resp));
  }

  static Future<int> createTicket(String categoryId, String message,
      List<FileInputInfo> attachments) async {
    var resp = await Server.post("/support/ticket/create",
        body: {"categoryId": categoryId, "message": message},
        files: _fromAttachment(attachments));
    return jsonDec(resp)["ticketId"];
  }

  static List<MultipartFile> _fromAttachment(List<FileInputInfo> attachments) {
    List<MultipartFile> files = [];
    var i = 0;
    for (final attachment in attachments) {
      i++;
      files.add(MultipartFile.fromBytes("attachment.$i", attachment.bytes(),
          filename: attachment.name()));
    }
    return files;
  }

  static Future<void> createTicketMessage(
      int ticketId, String message, List<FileInputInfo> attachments) async {
    await Server.post("/support/ticket/message",
        body: {"ticketId": ticketId.toString(), "message": message},
        files: _fromAttachment(attachments));
  }

  static Future<String> closeTicket(int ticketId) async {
    var resp = await Server.post("/support/ticket/close",
        query: {"id": ticketId.toString()});
    return jsonDec(resp)["message"];
  }

  static connectNotifier(int ticketId, Function() onNotified) {
    final channel = Server.createWebSocketChannel(
        "/ws-support/listenTicketUpdate?id=$ticketId");
    channel.stream.forEach((_) {
      onNotified();
    });
    return channel;
  }
}
