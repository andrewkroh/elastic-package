import json
from mitmproxy import http

def request(flow: http.HTTPFlow) -> None:
    # Check if the request is a POST to the _bulk endpoint and contains the required header
    if (flow.request.method == "POST" and
            flow.request.pretty_url.startswith("https://elasticsearch:9200/_bulk") and
            flow.request.headers.get("X-Elastic-Product-Origin") == "beats"):

        # Get the request body
        body = flow.request.content.decode("utf-8")

        # Split the body by lines (NDJSON format)
        lines = body.splitlines()

        # Iterate through the lines in pairs: first is control, second is the event
        for i in range(0, len(lines), 2):
            control_line = lines[i]
            event_line = lines[i + 1] if i + 1 < len(lines) else None

            # Parse the control line as JSON
            try:
                control_data = json.loads(control_line)
            except json.JSONDecodeError:
                continue  # Skip this pair if control line is not valid JSON

            # Check if the control message has the 'create' keyword
            if "create" in control_data and "_index" in control_data["create"]:
                data_stream = control_data["create"]["_index"];

                # Ignore agent logs and metrics.
                if data_stream.startswith("logs-elastic_agent") or data_stream.startswith("metrics-elastic_agent"):
                    continue

                # Append the event line to the output file if it exists
                if event_line:
                    with open(data_stream, "a") as f:
                        f.write(control_line + "\n")
                        f.write(event_line + "\n")

