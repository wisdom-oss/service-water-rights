import http
import logging
import typing

import amqp_rpc_client
import fastapi.security

import configuration
import enums
import exceptions
import models.amqp
import models.internal

# %% Clients needed for the security
__logger = logging.getLogger("security")


def is_authorized_user(
    scopes: fastapi.security.SecurityScopes,
    token_scopes: None | str = fastapi.Header(alias="X-Authenticated-Scope"),
    user_id: None | str = fastapi.Header(default=None, alias="X-Authenticated-Userid"),
) -> typing.Union[bool, str]:
    """
    Check if the user calling this service is authorized.

    This security dependency needs to be used as fast api dependency in the methods

    :param scopes: The scopes this used needs to have to access this service
    :type scopes: list
    :param token_scopes: The scopes the token is authorized for
    :type token_scopes: str
    :return: Status of the authorization
    :param user_id: The id of the user accessing the service
    :type user_id: str
    :rtype: bool
    :raises exceptions.APIException: The user is not authorized to access this service
    """
    _token_scopes = set([scope.strip() for scope in token_scopes.split(",")])
    print(_token_scopes)
    if "administrator" in token_scopes:
        return True if user_id is None else user_id
    required_scopes = set(scopes.scopes)
    if required_scopes.issubset(_token_scopes):
        return True if user_id is None else user_id
    else:
        raise exceptions.APIException(
            error_code="MISSING_PRIVILEGES",
            error_title="Missing Privileges",
            error_description="The account used to access this resource does not have the privileges to access "
            "this endpoint",
            http_status=http.HTTPStatus.FORBIDDEN,
        )
