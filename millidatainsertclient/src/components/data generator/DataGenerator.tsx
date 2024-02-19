'use client'
import React from 'react'
import { Button } from '../ui/button'
import { CsvDataArray, CsvValue, DataType } from './dataGenerator.types'
import { faker } from '@faker-js/faker'
import Papa from 'papaparse';
import axios, { AxiosResponse } from 'axios'


const ROW_SIZE = 10000

const DataGenerator = () => {


//   export const collectPagginatedDataInBatches = async ({
//     batchSize,
//     count,
//     requestFunction,
//     requestParams = {},
//     maxQtt,
// }: {
//     batchSize: number,
//     count: number,
//     requestFunction: (requestParams: GetAllRequestParamsType) => Promise<AxiosResponse<any, any>>
//     requestParams?: GetAllRequestParamsType
//     maxQtt?: number
// }) => {

//     const results: (AxiosResponse<any, any> | null | undefined)[] = [];
//     for (let batchIteration = 1; batchIteration <= pagesQtt; batchIteration += batchSize) {
//         const batch = Array.from({ length: batchSize }, (__, i) => i + batchIteration);
//         // eslint-disable-next-line no-await-in-loop
//         const batchResults = await Promise.all(batch
//             .map((iteration) => {
//                 if (iteration > pagesQtt) { return null; }
//                 return requestFunction({ ...requestParams, page: iteration, count });
//             }));
//         results.push(...batchResults);
//     }
//     return results;
// };

  const generateCSV = (headers: string[], values: CsvValue[]) => {
    // Create an array to hold the parsed data
    const data = [];

    // Push the headers as the first row
    data.push(headers);

    // Iterate over the values and separate them into the correct headers
    values.forEach((value) => {
        const rowData = headers.map((header) => value[header]);
        data.push(rowData);
    });

    // Use PapaParse to unparse the data into CSV format
    return Papa.unparse(data);
};

  const handleGenerateCSV = (csvDataArray: CsvDataArray) => {
        const { headers, values, fileName } = csvDataArray;
        if (!headers) return;
        const csvContent = generateCSV(headers, values);

        return new Blob([csvContent], { type: 'text/csv' });;

};

    const handlePostCsv = async (csvFile: Blob) => {
      const formData = new FormData();

      formData.append('csv',csvFile);
      axios.post('http://localhost:2030/posttodb',formData).then((res)=>{
        console.log(res);
      })
    }

    const handlePostMilliRequest = async () => {

      const data:any = [];
      for(let i=1;i<ROW_SIZE;i++){
          data.push({
            ProductId: i.toString(),
            Country: faker.location.country(),
            Location: faker.location.city(),
            Price: faker.number.int({min:100,max:100000}),
            Product: faker.word.noun(),
            RegisteredAt: faker.date.past().toISOString(),
            Shop: faker.company.name(),
            Type: faker.color.human()
          })
      }

      const batches = ROW_SIZE/1000;

      const postPromises: Promise<any>[] = [];


      for(let i=0;i<data.length;i++){
        const batchpromise = [];

        console.log(i,'i value');

        for(let j=0;j<50 && j<data.length;j++,i++){
          console.log(j,'j value');

          if(data[j+i])
            batchpromise.push(axios.post("http://localhost:2030/postmillirequests",data[j+i]));
        }

        await Promise.all(batchpromise).then((res)=>{
          console.log(res,'okokokok');
        })

      }

      // data.forEach((item)=>{
      //   postPromises.push(axios.post("http://localhost:2030/postmillirequests",item));
      // })

      // Promise.all(postPromises).then((res)=>{
      //   console.log(res,'okokokok');
      // })
    }

    const handleGenerateData = (isCsv: boolean) => {

        const headers = [
            'productId',
            'product',
            'shop',
            'location',
            'type',
            'price',
            'country',
            'registeredAt'
        ];
        const data:DataType[] = [];
        for(let i=1;i<ROW_SIZE;i++){
            data.push({
              productId: i.toString(),
              country: faker.location.country(),
              location: faker.location.city(),
              price: faker.number.int({min:100,max:100000}),
              product: faker.word.noun(),
              registeredAt: faker.date.past().toISOString(),
              shop: faker.company.name(),
              type: faker.color.human()
            })
        }
        const mycsv = handleGenerateCSV({
          fileName:'mydata',
          headers:headers,
          values:data as unknown as CsvValue[]
        })

        if(mycsv && isCsv){
          handlePostCsv(mycsv);
        }
    }
  return (
    <div>
      <Button onClick={()=>handleGenerateData(true)}>Generate CSV DATA</Button>
      <Button onClick={()=>handlePostMilliRequest()}>Generate Million Request</Button>

    </div>
  )
}

export default DataGenerator
