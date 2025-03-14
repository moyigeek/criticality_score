// This file is auto-generated by @hey-api/openapi-ts

import type { Options as ClientOptions, TDataShape, Client } from '@hey-api/client-fetch';
import type { GetHistoriesData, GetHistoriesResponse, GetQueryWithPaginationData, GetQueryWithPaginationResponse, GetRankingsData, GetRankingsResponse, GetResultsData, GetResultsResponse, GetResultsByScoreidData, GetResultsByScoreidResponse, GetSearchPackagesData, GetSearchPackagesResponse, PostUpdateGitlinkData, PostUpdateGitlinkResponse } from './types.gen';
import { client as _heyApiClient } from './client.gen';

export type Options<TData extends TDataShape = TDataShape, ThrowOnError extends boolean = boolean> = ClientOptions<TData, ThrowOnError> & {
    /**
     * You can provide a client instance returned by `createClient()` instead of
     * individual options. This might be also useful if you want to implement a
     * custom client.
     */
    client?: Client;
};

/**
 * Get score histories
 * Get score histories by git link
 */
export const getHistories = <ThrowOnError extends boolean = false>(options: Options<GetHistoriesData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).get<GetHistoriesResponse, unknown, ThrowOnError>({
        url: '/histories',
        ...options
    });
};

/**
 * Query with pagination
 * Query the database with pagination support
 */
export const getQueryWithPagination = <ThrowOnError extends boolean = false>(options: Options<GetQueryWithPaginationData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).get<GetQueryWithPaginationResponse, unknown, ThrowOnError>({
        url: '/query-with-pagination',
        ...options
    });
};

/**
 * Get ranking results
 * Get ranking results, optionally including all details
 */
export const getRankings = <ThrowOnError extends boolean = false>(options?: Options<GetRankingsData, ThrowOnError>) => {
    return (options?.client ?? _heyApiClient).get<GetRankingsResponse, unknown, ThrowOnError>({
        url: '/rankings',
        ...options
    });
};

/**
 * Search score results by git link
 * Search score results by git link
 * NOTE: All details are ignored, should use /results/:scoreid to get details
 * NOTE: Maxium take count is 1000
 */
export const getResults = <ThrowOnError extends boolean = false>(options: Options<GetResultsData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).get<GetResultsResponse, unknown, ThrowOnError>({
        url: '/results',
        ...options
    });
};

/**
 * Get score results
 * Get score results, including all details by scoreid
 */
export const getResultsByScoreid = <ThrowOnError extends boolean = false>(options: Options<GetResultsByScoreidData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).get<GetResultsByScoreidResponse, unknown, ThrowOnError>({
        url: '/results/{scoreid}',
        ...options
    });
};

/**
 * Search packages
 * Search for packages in the specified table that match the search query
 */
export const getSearchPackages = <ThrowOnError extends boolean = false>(options: Options<GetSearchPackagesData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).get<GetSearchPackagesResponse, unknown, ThrowOnError>({
        url: '/search-packages',
        ...options
    });
};

/**
 * Update git link
 * Update the git link for a specified package
 */
export const postUpdateGitlink = <ThrowOnError extends boolean = false>(options: Options<PostUpdateGitlinkData, ThrowOnError>) => {
    return (options.client ?? _heyApiClient).post<PostUpdateGitlinkResponse, unknown, ThrowOnError>({
        url: '/update-gitlink',
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers
        }
    });
};