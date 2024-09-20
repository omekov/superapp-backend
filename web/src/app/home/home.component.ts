import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { HttpService, IMark, IModel, ISpecification, IVolume } from '../http.service';
import { FormsModule } from '@angular/forms';
import { CustomCurrencyPipe } from '../custom-currency.pipe';
import { Router } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule, CustomCurrencyPipe],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css',
})
export class HomeComponent implements OnInit {
  marks: IMark[] = []
  mark: string = "";
  models: IModel[] = []
  model: string = "";
  volumes: IVolume[] = [];
  volume: number = 0;
  specifications: ISpecification[] = [];
  specification: ISpecification = { Year: 0, Amount: 0 };
  amount: number = 0;
  amountKZT: number = 0;
  USDToKZT: number = 480;
  year: number = 2000;
  mrp: number = 3692;
  customsCollection: number = 0;
  customsDuty: number = 0;
  vat: number = 0;
  deliveredAmount: number = 2100;
  constructor(private httpService: HttpService,
    private route: ActivatedRoute,
    private router: Router) { }

  ngOnInit() {
    this.customsCollection = this.mrp * 6
    this.deliveredAmount = this.deliveredAmount * this.USDToKZT

    this.route.queryParams.subscribe(params => {
      this.mark = params['mark'] || null;
      this.model = params['model'] || null;
      this.volume = params['volume'] || null;
      this.year = params['year'] || null;
    });

    this.httpService.getMarks().subscribe((data: IMark[]) => {
      this.marks = data
    });
    if (this.mark !== "" && this.model === "" && this.volume === 0 && this.year === 0) {
      this.httpService.getModels(this.mark).subscribe((data: IModel[]) => {
        this.models = data
      });
    } else if (this.mark !== "" && this.model !== "" && this.volume === 0 && this.year === 0) {
      this.httpService.getModels(this.mark).subscribe((data: IModel[]) => {
        this.models = data
      });
      this.httpService.getVolumes(this.mark, this.model).subscribe((data: IVolume[]) => {
        this.volumes = data
      });
    } else if (this.mark !== "" && this.model !== "" && this.volume !== 0 && this.year !== 0) {
      this.httpService.getModels(this.mark).subscribe((data: IModel[]) => {
        this.models = data
      });
      this.httpService.getVolumes(this.mark, this.model).subscribe((data: IVolume[]) => {
        this.volumes = data
      });
      this.httpService.getSpecifications(this.mark, this.model, this.volume).subscribe((data: ISpecification[]) => {
        this.specifications = data

        this.specifications.map(data => {
          if (data.Year == this.year) {
            this.amount = data.Amount
            return
          }
        });
      });
    }
  }

  onChangeMark(event: any) {
    this.mark = event.target.value
    this.httpService.getModels(this.mark).subscribe((data: IModel[]) => {
      this.models = data
      this.router.navigate([], {
        queryParams: {
          mark: this.mark,
        },
        queryParamsHandling: 'merge', // Чтобы сохранить другие параметры в URL
      });
    });
  }

  onChangeModel(event: any) {
    this.model = event.target.value
    this.httpService.getVolumes(this.mark, this.model).subscribe((data: IVolume[]) => {
      this.volumes = data

      this.router.navigate([], {
        queryParams: {
          mark: this.mark,
          model: this.model,
        },
        queryParamsHandling: 'merge', // Чтобы сохранить другие параметры в URL
      });
    });
  }


  onChangeVolume(event: any) {
    this.volume = event.target.value
    this.httpService.getSpecifications(this.mark, this.model, this.volume).subscribe((data: ISpecification[]) => {
      this.specifications = data

      this.router.navigate([], {
        queryParams: {
          mark: this.mark,
          model: this.model,
          volume: this.volume,
        },
        queryParamsHandling: 'merge', // Чтобы сохранить другие параметры в URL
      });
    });
  }

  onChangeSpecification(event: any) {
    this.amount = event.target.value
    this.specifications.map(data => {
      this.year = data.Year
    })
    this.router.navigate([], {
      queryParams: {
        mark: this.mark,
        model: this.model,
        volume: this.volume,
        year: this.year
      },
      queryParamsHandling: 'merge', // Чтобы сохранить другие параметры в URL
    });

    this.amountKZT = this.amount * this.USDToKZT
    this.customsDuty = (this.amountKZT * 15) / 100
    this.vat = ((this.amountKZT + this.customsDuty + this.customsCollection) * 12) / 100
  }
}
