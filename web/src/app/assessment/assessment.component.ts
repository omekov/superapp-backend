import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Location } from '@angular/common';
@Component({
  selector: 'app-assessment',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './assessment.component.html',
  styleUrl: './assessment.component.css'
})
export class AssessmentComponent implements OnInit {

  constructor(private location: Location, ) {}

  ngOnInit(): void {
      
  }

  mrp: number = 3692;
  usd: number = 480.55;
  usd_update_date = new Date();
  goBack() {
    this.location.back();
  }



}
