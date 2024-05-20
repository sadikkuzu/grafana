import { BaseQueryFn } from '@reduxjs/toolkit/query/react';
import { lastValueFrom } from 'rxjs';

import { getBackendSrv, isFetchError } from '@grafana/runtime/src/services/backendSrv';

import { DataQuerySpecResponse } from './types';

/**
 * @alpha
 */
export const API_VERSION = 'peakq.grafana.app/v0alpha1';

/**
 * @alpha
 */
export enum QueryTemplateKinds {
  QueryTemplate = 'QueryTemplate',
}

/**
 * Query Library is an experimental feature. API (including the URL path) will likely change.
 *
 * @alpha
 */
export const BASE_URL = `/apis/${API_VERSION}/namespaces/default/querytemplates/`;

/**
 *
 * @alpha
 */
export type QueryTemplateRequestOptions = {
  method?: 'POST';
  data?: any;
};

/**
 * TODO: similar code is duplicated in many places. To be unified in #86960
 */
export const baseQuery: BaseQueryFn<QueryTemplateRequestOptions, DataQuerySpecResponse, Error> = async (
  requestOptions
) => {
  try {
    const responseObservable = getBackendSrv().fetch<DataQuerySpecResponse>({
      url: BASE_URL,
      showErrorAlert: true,
      method: requestOptions.method || 'GET',
      data: requestOptions.data,
    });
    return await lastValueFrom(responseObservable);
  } catch (error) {
    if (isFetchError(error)) {
      return { error: new Error(error.data.message) };
    } else if (error instanceof Error) {
      return { error };
    } else {
      return { error: new Error('Unknown error') };
    }
  }
};
