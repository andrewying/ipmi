/*
Source: https://github.com/pion/example-webrtc-applications

MIT License

Copyright (c) 2018

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */

#include "gst.h"

#include <gst/app/gstappsrc.h>

typedef struct SampleHandlerUserData {
  int pipelineId;
} SampleHandlerUserData;

GMainLoop *gstreamer_send_main_loop = NULL;
void gstreamer_send_start_mainloop(void) {
  gstreamer_send_main_loop = g_main_loop_new(NULL, FALSE);

  g_main_loop_run(gstreamer_send_main_loop);
}

static gboolean gstreamer_send_bus_call(GstBus *bus, GstMessage *msg, gpointer data) {
  switch (GST_MESSAGE_TYPE(msg)) {

  case GST_MESSAGE_EOS:
    g_print("End of stream\n");
    exit(1);
    break;

  case GST_MESSAGE_ERROR: {
    gchar *debug;
    GError *error;

    gst_message_parse_error(msg, &error, &debug);
    g_free(debug);

    g_printerr("Error: %s\n", error->message);
    g_error_free(error);
    exit(1);
  }
  default:
    break;
  }

  return TRUE;
}

GstFlowReturn gstreamer_send_new_sample_handler(GstElement *object, gpointer user_data) {
  GstSample *sample = NULL;
  GstBuffer *buffer = NULL;
  gpointer copy = NULL;
  gsize copy_size = 0;
  SampleHandlerUserData *s = (SampleHandlerUserData *)user_data;

  g_signal_emit_by_name (object, "pull-sample", &sample);
  if (sample) {
    buffer = gst_sample_get_buffer(sample);
    if (buffer) {
      gst_buffer_extract_dup(buffer, 0, gst_buffer_get_size(buffer), &copy, &copy_size);
      goHandlePipelineBuffer(copy, copy_size, GST_BUFFER_DURATION(buffer), s->pipelineId);
    }
    gst_sample_unref (sample);
  }

  return GST_FLOW_OK;
}

GstElement *gstreamer_send_create_pipeline(char *pipeline) {
  gst_init(NULL, NULL);
  GError *error = NULL;
  return gst_parse_launch(pipeline, &error);
}

void gstreamer_send_start_pipeline(GstElement *pipeline, int pipelineId) {
  SampleHandlerUserData *s = calloc(1, sizeof(SampleHandlerUserData));
  s->pipelineId = pipelineId;

  GstBus *bus = gst_pipeline_get_bus(GST_PIPELINE(pipeline));
  gst_bus_add_watch(bus, gstreamer_send_bus_call, NULL);
  gst_object_unref(bus);

  GstElement *appsink = gst_bin_get_by_name(GST_BIN(pipeline), "appsink");
  g_object_set(appsink, "emit-signals", TRUE, NULL);
  g_signal_connect(appsink, "new-sample", G_CALLBACK(gstreamer_send_new_sample_handler), s);
  gst_object_unref(appsink);

  gst_element_set_state(pipeline, GST_STATE_PLAYING);
}

void gstreamer_send_stop_pipeline(GstElement *pipeline) {
  gst_element_set_state(pipeline, GST_STATE_NULL);
}
