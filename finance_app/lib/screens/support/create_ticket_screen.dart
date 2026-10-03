import 'package:finance_app/model/support.dart';
import 'package:finance_app/service/support_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/file_input.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/select/select_box.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';

import '../../utils/screens.dart';
import 'help_ticket_screen.dart';

class CreateTicketScreen extends StatefulWidget {
  const CreateTicketScreen(
      {Key? key, this.initialMessage, this.initialProblemCategory})
      : super(key: key);

  final String? initialMessage;
  final String? initialProblemCategory;

  @override
  State<StatefulWidget> createState() => _CreateTicketState();
}

class _CreateTicketState extends State<CreateTicketScreen> {
  final messageController = TextEditingController();
  final problemCategoryController = SelectController<String?>(null);
  final attachmentController = FileInputController();

  List<SupportTicketCategoryModel>? problemCategories;

  @override
  void initState() {
    super.initState();

    messageController.text = widget.initialMessage ?? "";
    if (widget.initialProblemCategory != null) {
      problemCategoryController.setValue(widget.initialProblemCategory);
    }

    SupportService.getTicketCategories().then((v) {
      setState(() {
        problemCategories = v;
      });
    });
  }

  void _onSubmit() async {
    String message = messageController.text;
    String? categoryId = problemCategoryController.value();
    if (categoryId == null) {
      Alert.message("Jenis masalah harus dipilih");
      return;
    }
    if (message.isEmpty) {
      Alert.message("Pesan harus diisi");
      return;
    }

    Alert.withLoading((done) async {
      var ticketId = await SupportService.createTicket(
          categoryId, message, attachmentController.value);
      done();
      Screens.replace(HelpTicketScreen(ticketId: ticketId));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Buat Tiket Bantuan"),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SelectBox(
                  controller: problemCategoryController,
                  enabled: problemCategories != null,
                  items: (problemCategories ?? [])
                      .map((e) => SelectItem(e.id, e.label))
                      .toList(),
                  label: "Jenis Masalah"),
              const SizedBox(height: 16),
              Input(
                controller: messageController,
                label: "Pesan",
                keyboardType: TextInputType.multiline,
                maxLines: 12,
                minLines: 6,
                padding:
                    const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
              ),
              const SizedBox(height: 16),
              FileInput(
                label: "Lampiran",
                controller: attachmentController,
              ),
              const SizedBox(height: 16),
              Button(
                  onPressed: _onSubmit,
                  width: double.infinity,
                  child: const Text("Kirim"))
            ],
          ),
        ),
      ),
    );
  }
}
