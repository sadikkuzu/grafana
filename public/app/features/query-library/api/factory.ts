import { createApi } from '@reduxjs/toolkit/query/react';

import { AddQueryTemplateCommand, QueryTemplate } from '../types';

import { convertAddQueryTemplateCommandToDataQuerySpec, convertDataQueryResponseToQueryTemplates } from './mappers';
import { baseQuery } from './query';

export const queryLibraryApi = createApi({
  baseQuery,
  endpoints: (builder) => ({
    allQueryTemplates: builder.query<QueryTemplate[], void>({
      query: () => ({}),
      transformResponse: convertDataQueryResponseToQueryTemplates,
    }),
    addQueryTemplate: builder.mutation<QueryTemplate, AddQueryTemplateCommand>({
      query: (addQueryTemplateCommand) => ({
        method: 'POST',
        data: convertAddQueryTemplateCommandToDataQuerySpec(addQueryTemplateCommand),
      }),
    }),
  }),
  reducerPath: 'queryLibrary',
});
