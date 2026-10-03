import 'dart:typed_data';

import 'package:file_picker/file_picker.dart';
import 'package:finance_app/screens/utility/image_viewer_screen.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';

class FileInput extends StatefulWidget {
  const FileInput({Key? key, this.label, this.controller}) : super(key: key);

  final String? label;
  final FileInputController? controller;

  @override
  State<StatefulWidget> createState() => _FileInputState();
}

class FileInputInfo {
  final PlatformFile _file;

  const FileInputInfo(this._file);

  String name() {
    return _file.path!.split("/").last;
  }

  String ext() {
    return _file.path!.split(".").last;
  }

  String path() {
    return _file.path!;
  }

  ImageProvider imageProvider() {
    return MemoryImage(_file.bytes!);
  }

  Uint8List bytes() {
    return _file.bytes!;
  }

  bool isImage() {
    return Meta.isImageExt(name());
  }

  IconData icon() {
    if (isImage()) {
      return Icons.image;
    }
    return Icons.question_mark;
  }
}

class FileInputController extends ValueNotifier<List<FileInputInfo>> {
  FileInputController({List<FileInputInfo>? files}) : super(files ?? []);

  void removeFile(FileInputInfo file) {
    value.remove(file);
    notifyListeners();
  }

  void clean() {
    value.clear();
  }

  void addFiles(List<FileInputInfo> files) {
    for (final file in files) {
      if (value.where((e) => e.path() == file.path()).isEmpty) {
        value.add(file);
      }
    }
    notifyListeners();
  }
}

class _FileInputState extends State<FileInput> {
  late FileInputController controller;

  @override
  void initState() {
    super.initState();
    controller = widget.controller ?? FileInputController();
    controller.addListener(() {
      if (mounted) {
        setState(() {});
      }
    });
  }

  void _onPickFile() async {
    FilePickerResult? result = await FilePicker.platform
        .pickFiles(allowMultiple: true, withData: true);

    if (result != null) {
      List<FileInputInfo> files =
          result.files.map((f) => FileInputInfo(f)).toList();
      controller.addFiles(files);
    }
  }

  @override
  void dispose() {
    super.dispose();
    controller.clean();
  }

  Widget _buildFile(FileInputInfo fileInfo) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: ColorHelper.inputBg(context),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Row(
        children: [
          Icon(fileInfo.icon()),
          const SizedBox(width: 16),
          Expanded(
              child: GestureDetector(
            onTap: () {
              if (!fileInfo.isImage()) {
                return;
              }
              Screens.to(ImageViewerScreen(
                  title: fileInfo.name(), provider: fileInfo.imageProvider()));
            },
            child: Text(fileInfo.name()),
          )),
          const SizedBox(width: 16),
          GestureDetector(
            onTap: () {
              controller.removeFile(fileInfo);
            },
            child: const Icon(Icons.close, size: 20),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (widget.label != null) ...[
          Text(widget.label!, style: Input.labelTextStyle(context)),
          const SizedBox(height: 8)
        ],
        ...controller.value.map((e) => _buildFile(e)).toList(),
        GestureDetector(
          onTap: _onPickFile,
          child: Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: ColorHelper.inputBg(context),
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Icon(
              Icons.upload,
              size: 30,
            ),
          ),
        )
      ],
    );
  }
}
