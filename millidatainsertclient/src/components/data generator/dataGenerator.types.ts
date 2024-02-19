export interface DataType{
    productId: string,
    product: string,
    shop: string,
    price: number
    location: string
    type: string
    country: string
    registeredAt: string
}
export interface CsvValue { [header: string]: string | number | undefined }


export interface CsvDataArray {
    headers: string[] | null
    values: CsvValue[]
    fileName: string
}