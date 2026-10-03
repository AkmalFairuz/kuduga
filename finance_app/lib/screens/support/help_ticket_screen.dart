import 'package:cached_network_image/cached_network_image.dart';
import 'package:finance_app/model/support.dart';
import 'package:finance_app/screens/utility/image_viewer_screen.dart';
import 'package:finance_app/service/support_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/bottom_buttons.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/modal_bottom_sheet.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/file_input.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/skeleton/skeleton.dart';
import 'package:finance_app/widgets/text/modal_bottom_title.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class HelpTicketScreen extends StatefulWidget {
  const HelpTicketScreen({Key? key, required this.ticketId}) : super(key: key);

  final int ticketId;

  @override
  State<StatefulWidget> createState() => _HelpTicketState();
}

class _HelpTicketState extends State<HelpTicketScreen> {
  final messageController = TextEditingController();
  final scrollController = ScrollController();
  late WebSocketChannel channel;
  bool _closing = false;

  SupportTicketModel? _ticket;

  @override
  void initState() {
    super.initState();
    loadTicket();
  }

  @override
  void dispose() {
    _closing = true;
    channel.sink.close();
    super.dispose();
  }

  Future<void> loadTicket() async {
    if (!mounted) {
      return;
    }
    try {
      var resp = await SupportService.getDetailedTicket(widget.ticketId);
      _ticket = resp;
      channel = SupportService.connectNotifier(widget.ticketId, _onRefresh);

      channel.sink.done.then((_) async {
        if (_closing) {
          return;
        }
        await loadTicket();
      });
      setState(() {});
    } catch (e) {
      Future.delayed(const Duration(seconds: 3)).then((value) {
        loadTicket();
      });
    }
  }

  Widget _buildMessages() {
    if (_ticket == null) {
      return const Column(
        children: [
          Expanded(child: SizedBox()),
          _TicketMessage(),
          Line(),
          _TicketMessage(),
          Line(),
          _TicketMessage(),
          Line(),
          _TicketMessage(),
        ],
      );
    }
    int msgLen = _ticket!.messages!.length;
    return RefreshIndicator(
        onRefresh: loadTicket,
        child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            reverse: true,
            controller: scrollController,
            itemBuilder: (context, index) {
              var msg =
                  _ticket!.messages![_ticket!.messages!.length - index - 1];
              return _TicketMessage(key: Key(index.toString()), model: msg);
            },
            separatorBuilder: (_, __) => const Line(),
            itemCount: msgLen));
  }

  Future<bool> _onSubmit(String message, List<FileInputInfo> files) async {
    if (message.isEmpty) {
      return false;
    }
    await SupportService.createTicketMessage(widget.ticketId, message, files);
    return true;
  }

  void _onSubmitAttachment(String message, List<FileInputInfo> files) async {
    await Alert.withLoading((done) async {
      var r = await _onSubmit(message, files);
      done();
      if (r) {
        Screens.back();
      }
    });
  }

  void _onAttachment() {
    ModalBottomSheet.show(context, (context) {
      return _MessageAttachmentForm(onSubmit: _onSubmitAttachment);
    }, showDragHelper: false);
  }

  void _onRefresh() async {
    final lastId = _ticket!.messages!.last.id;
    final ticket = await SupportService.getDetailedTicket(_ticket!.id,
        messageAfterId: lastId);
    _ticket!.status = ticket.status;
    if (ticket.messages!.isNotEmpty) {
      // safe add messages (prevent duplicate)
      _ticket!.messages!.addAll(ticket.messages!
          .where((e) => !_ticket!.messages!.contains(e))
          .toList());

      Future.delayed(const Duration(milliseconds: 100)).then((_) {
        if (!mounted) return;
        scrollController.animateTo(0,
            duration: const Duration(milliseconds: 400),
            curve: Curves.fastOutSlowIn);
      });
    }
    if (mounted) {
      setState(() {});
    }
  }

  void _onCloseTicket() async {
    var message = await SupportService.closeTicket(_ticket!.id);
    Alert.message(message);
    setState(() {
      _ticket!.status = 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: AppBarTitle("Tiket #${widget.ticketId}"),
        actions: _ticket != null && _ticket!.status == 0
            ? [
                IconButton(
                  onPressed: _onRefresh,
                  icon: const Icon(Icons.refresh),
                  tooltip: "Refresh",
                ),
                IconButton(
                    onPressed: () {
                      BottomButtons.show(context, [
                        BottomButtonItem(
                            text: "Tutup Tiket", onTap: _onCloseTicket),
                      ]);
                    },
                    icon: const Icon(Icons.more_vert)),
              ]
            : [],
      ),
      body: Column(
        children: [
          Expanded(child: _buildMessages()),
          const Line(),
          if (_ticket != null && _ticket!.status != 0)
            const Padding(
                padding: EdgeInsets.all(8), child: Text("Tiket sudah ditutup")),
          if (_ticket != null && _ticket!.status == 0)
            Container(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Expanded(
                      child: Input(
                    controller: messageController,
                    padding: const EdgeInsets.all(14),
                    maxLines: 2,
                    minLines: 1,
                    hintText: "Pesan",
                    keyboardType: TextInputType.multiline,
                    suffixIcon: IntrinsicWidth(
                      child: Row(
                        children: [
                          const SizedBox(width: 16),
                          GestureDetector(
                              onTap: _onAttachment,
                              child: const Icon(Icons.attachment)),
                          const SizedBox(width: 20),
                          GestureDetector(
                              onTap: () {
                                _onSubmit(messageController.text, []);
                                messageController.clear();
                              },
                              child: const Icon(Icons.send)),
                          const SizedBox(width: 16),
                        ],
                      ),
                    ),
                  )),
                ],
              ),
            )
        ],
      ),
    );
  }
}

class _TicketMessage extends StatelessWidget {
  const _TicketMessage({Key? key, this.model}) : super(key: key);

  final SupportTicketMessageModel? model;

  Widget _buildAttachment(BuildContext context, String url) {
    return IntrinsicWidth(
        child: GestureDetector(
      onTap: () {
        if (Meta.isImageExt(url)) {
          Screens.to(ImageViewerScreen(
              title: Meta.fileWithoutPath(url),
              provider: CachedNetworkImageProvider(url)));
        }
      },
      child: Container(
        decoration: BoxDecoration(
            color: ColorHelper.bg200(context),
            borderRadius: BorderRadius.circular(8)),
        padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 8),
        child: Row(
          children: [
            Icon(Meta.isImageExt(url) ? Icons.image : Icons.question_mark),
            const SizedBox(width: 8),
            Flexible(
              child: Text(
                Meta.fileWithoutPath(url),
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ],
        ),
      ),
    ));
  }

  Widget _buildAttachments(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 8),
        ...model!.attachments.map((e) => _buildAttachment(context, e)).toList()
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
        color: model?.role == 1
            ? ColorHelper.theme(context,
                dark: Meta.color[900], light: Meta.color[50])
            : null,
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                if (model != null && model!.role == 1) ...[
                  Icon(
                    Icons.support_agent,
                    size: 18,
                    color: Meta.color,
                  ),
                  const SizedBox(width: 8),
                ],
                model != null
                    ? Text(
                        model!.role == 0 ? "Anda" : model!.author,
                        style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: model!.role == 1 ? Meta.color : null),
                      )
                    : const Skeleton(width: 90, height: 15),
                const SizedBox(width: 12),
                model != null
                    ? Text(
                        Meta.formatUnixDate(model!.createdAt),
                        style: TextStyle(color: Colors.grey[500], fontSize: 11),
                      )
                    : const Skeleton(width: 60, height: 12),
              ],
            ),
            SizedBox(height: model == null ? 10 : 4),
            model != null
                ? Text(model!.message)
                : const Skeleton(width: double.infinity, height: 50),
            if (model != null && model!.attachments.isNotEmpty)
              _buildAttachments(context)
          ],
        ));
  }
}

class _MessageAttachmentForm extends StatefulWidget {
  const _MessageAttachmentForm({Key? key, required this.onSubmit})
      : super(key: key);

  final Function(String, List<FileInputInfo>) onSubmit;

  @override
  State<StatefulWidget> createState() => _MessageAttachmentFormState();
}

class _MessageAttachmentFormState extends State<_MessageAttachmentForm> {
  final FileInputController attachmentController = FileInputController();
  final TextEditingController messageController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const ModalBottomTitle(
          "Kirim pesan dengan lampiran",
          fontSize: 18,
        ),
        const SizedBox(height: 16),
        FileInput(
          controller: attachmentController,
        ),
        const SizedBox(height: 16),
        Input(hintText: "Pesan", controller: messageController),
        const SizedBox(height: 16),
        Button(
            onPressed: () {
              widget.onSubmit(
                  messageController.text, attachmentController.value);
            },
            width: double.infinity,
            child: const Text("Kirim"))
      ],
    );
  }
}
