import requests
import json
import os
from behave import given, when, then


@given("protocol {protocol}")
def protocol(context, protocol):
    setattr(context, "protocol", protocol)


@given("host {host}")
def host(context, host):
    setattr(context, "host", host)


@given("port {port}")
def port(context, port):
    setattr(context, "port", port)


@given("path {path}")
def path(context, path):
    setattr(context, "path", path)


@given("header {key} as {value}")
def header_key_as_value(context, key, value):
    if not hasattr(context, "request_headers"):
        context.request_headers = {}

    context.request_headers[key] = value
    print("header: " + context.request_headers[key])


def get_value(context, attribute, default):
    env_name = "TESTS_" + attribute.upper()

    if hasattr(context, attribute):
        return getattr(context, attribute)
    else:
        if env_name in os.environ:
            return os.environ.get(env_name)
        else:
            return default


@given("request {request}")
def request_body(context, request):
    dir = os.path.dirname(__file__)
    try:
        file_name = os.path.join(dir, "../requests", request)
        print(file_name)
        file = open(file_name, "r")
        body = json.loads(file.read())

        context.request_body = body
        print(context.request_body)
    finally:
        file.close()


def get_request_url(context):
    protocol = get_value(context, "protocol", "http")
    host = get_value(context, "host", "localhost")
    port = get_value(context, "port", "8080")
    path = context.path
    url = f"{protocol}://{host}:{port}{path}"
    print(url)
    return url


@when("method {verb}")
def method_verb(context, verb):
    context.request_method = verb
    print("method: " + context.request_method)

    if not hasattr(context, "request_headers"):
        context.request_headers = {}

    if not hasattr(context, "request_body"):
        context.request_body = None

    if context.request_method == "GET":
        response = requests.get(
            get_request_url(context),
            headers=context.request_headers,
            json=context.request_body,
        )

    if context.request_method == "POST":
        response = requests.post(
            get_request_url(context),
            headers=context.request_headers,
            json=context.request_body,
        )

    context.response_body = response.text
    context.response_code = response.status_code
    print("actual: " + str(context.response_code))


@then("status {code}")
def status_code(context, code):
    print("expected: " + str(code))
    assert str(context.response_code) == str(code)
