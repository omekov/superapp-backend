// custom-currency.pipe.ts
import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'kzt',
  standalone: true
})
export class CustomCurrencyPipe implements PipeTransform {

  transform(value: number): string {
    if (!value) return '';

    const formattedValue = value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');

    return `${formattedValue} тг.`;
  }
}