#!/usr/bin/env python

import argparse
import os
import requests

# Disable urllib3 warnings:
# https://github.com/shazow/urllib3/issues/497
import requests.packages.urllib3
requests.packages.urllib3.disable_warnings()


def _assert_digest_not_found(base_url, cert_path, digest):
    url = '{}/messages/{}'.format(base_url, digest)
    r = requests.get(url=url, verify=cert_path)
    assert r.status_code == 404, r.text
    print('{} correctly not found'.format(url))


def _assert_digest_found(base_url, cert_path, digest, expected_message):
    url = '{}/messages/{}'.format(base_url, digest)
    r = requests.get(url=url, verify=cert_path)
    assert r.status_code == 200, r.text
    resp_json = r.json()
    assert 'message' in resp_json, 'message not in json: {}'.format(resp_json)
    assert resp_json['message'] == expected_message, r.text
    print('{} correctly found'.format(url))


def post_message(base_url, cert_path, message, expected_digest=None):
    url = '{}/messages'.format(base_url)
    r = requests.post(url=url, json={'message': message}, verify=cert_path)
    assert r.status_code == 201, r.text
    resp_json = r.json()
    assert 'digest' in resp_json, 'digest not in json: {}'.format(resp_json)
    returned_digest = resp_json['digest']
    if expected_digest:
        assert returned_digest == expected_digest, r.text
    print('{} POSTed successfully'.format(url))
    return returned_digest


if __name__ == '__main__':

    # Parse arguments
    parser = argparse.ArgumentParser(
        description='Script to test your micro-service is correctly configured',
    )
    parser.add_argument(
        '-d', '--domain',
        default='localhost',
        help='API domain (defaults to localhost)',
        required=False,
    )
    parser.add_argument(
        '-p', '--port',
        default='5000',
        help='API port (defaults to 5000)',
        required=False,
    )
    parser.add_argument(
        '-c', '--cert-path',
        default='localhost.crt',
        help='Path to self signed certificate (defaults to localhost.crt)',
        required=False,
    )
    args = parser.parse_args()

    if not os.path.exists(args.cert_path):
        raise Exception('No Certificate File Found at {}'.format(args.cert_path))

    BASE_URL = 'https://{}:{}'.format(args.domain, args.port)

    # GET a digest for a message that doesn't exist
    _assert_digest_not_found(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        digest='a' * 64,
    )

    # POST the message 'foo'
    DIGEST = '2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae'
    digest_returned = post_message(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        message='foo',
        expected_digest=DIGEST,
    )

    # GET the message 'foo' back
    _assert_digest_found(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        digest=DIGEST,
        expected_message='foo',
    )

    # POST the message 'bar'
    bar_digest_returned = post_message(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        message='bar',
    )

    # GET the message 'bar' back
    _assert_digest_found(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        digest=bar_digest_returned,
        expected_message='bar',
    )

    # GET the message 'foo' again
    _assert_digest_found(
        base_url=BASE_URL,
        cert_path=args.cert_path,
        digest=DIGEST,
        expected_message='foo',
    )

    print('')
    print('*' * 75)
    print('All tests passed!')
    print('*' * 75)

